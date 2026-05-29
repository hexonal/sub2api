package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

const entroxCLISessionTTL = 10 * time.Minute

type entroxCLISession struct {
	ID        string
	PollToken string
	APIKey    string
	APIKeyID  int64
	Created   bool
	Approving bool
	UserID    int64
	CreatedAt time.Time
	ExpiresAt time.Time
}

type approveEntroxCLIRequest struct {
	SessionID string `json:"session_id" binding:"required"`
	APIKeyID  *int64 `json:"api_key_id"`
	CreateNew bool   `json:"create_new"`
}

type approveEntroxCLIResponse struct {
	Status   string `json:"status"`
	APIKeyID int64  `json:"api_key_id,omitempty"`
	Created  bool   `json:"created"`
}

func (h *AuthHandler) EntroxCLIScript(c *gin.Context) {
	c.Data(http.StatusOK, "text/x-shellscript; charset=utf-8", []byte(entroxCLILoginScript))
}

func (h *AuthHandler) StartEntroxCLIAuth(c *gin.Context) {
	sessionID, err := randomHex(16)
	if err != nil {
		response.InternalError(c, "Failed to create auth session")
		return
	}
	pollToken, err := randomHex(32)
	if err != nil {
		response.InternalError(c, "Failed to create auth session")
		return
	}

	now := time.Now()
	h.entroxCLIMu.Lock()
	h.cleanupExpiredEntroxCLISessionsLocked(now)
	h.entroxCLISessions[sessionID] = &entroxCLISession{
		ID:        sessionID,
		PollToken: pollToken,
		CreatedAt: now,
		ExpiresAt: now.Add(entroxCLISessionTTL),
	}
	h.entroxCLIMu.Unlock()

	origin := authRequestOrigin(c)
	c.JSON(http.StatusOK, gin.H{
		"session_id":    sessionID,
		"poll_token":    pollToken,
		"authorize_url": origin + "/entrox/connect?session_id=" + sessionID,
		"expires_in":    int(entroxCLISessionTTL.Seconds()),
	})
}

func (h *AuthHandler) PollEntroxCLIAuth(c *gin.Context) {
	sessionID := strings.TrimSpace(c.Query("session_id"))
	pollToken := strings.TrimSpace(c.Query("poll_token"))
	if sessionID == "" || pollToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "invalid_request"})
		return
	}

	now := time.Now()
	h.entroxCLIMu.Lock()
	defer h.entroxCLIMu.Unlock()
	h.cleanupExpiredEntroxCLISessionsLocked(now)

	session := h.entroxCLISessions[sessionID]
	if session == nil || session.PollToken != pollToken {
		c.JSON(http.StatusNotFound, gin.H{"status": "not_found"})
		return
	}
	if now.After(session.ExpiresAt) {
		delete(h.entroxCLISessions, sessionID)
		c.JSON(http.StatusGone, gin.H{"status": "expired"})
		return
	}
	if session.APIKey == "" {
		c.JSON(http.StatusAccepted, gin.H{"status": "pending"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "approved",
		"api_key": session.APIKey,
	})
}

func (h *AuthHandler) ApproveEntroxCLIAuth(c *gin.Context) {
	if h.apiKeyService == nil {
		response.InternalError(c, "API key service is not available")
		return
	}

	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req approveEntroxCLIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	if err := validateApproveEntroxCLIRequest(req); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	sessionID := strings.TrimSpace(req.SessionID)

	now := time.Now()
	h.entroxCLIMu.Lock()
	h.cleanupExpiredEntroxCLISessionsLocked(now)
	session := h.entroxCLISessions[sessionID]
	if session == nil {
		h.entroxCLIMu.Unlock()
		response.NotFound(c, "Auth session not found")
		return
	}
	if now.After(session.ExpiresAt) {
		delete(h.entroxCLISessions, sessionID)
		h.entroxCLIMu.Unlock()
		response.Error(c, http.StatusGone, "Auth session expired")
		return
	}
	if session.APIKey != "" {
		resp := approveEntroxCLIResponse{
			Status:   "approved",
			APIKeyID: session.APIKeyID,
			Created:  session.Created,
		}
		h.entroxCLIMu.Unlock()
		response.Success(c, resp)
		return
	}
	if session.Approving {
		h.entroxCLIMu.Unlock()
		response.Success(c, gin.H{"status": "pending"})
		return
	}
	session.Approving = true
	h.entroxCLIMu.Unlock()

	key, created, err := h.resolveEntroxCLIAPIKey(c.Request.Context(), subject.UserID, req, now)
	if err != nil {
		h.entroxCLIMu.Lock()
		if session := h.entroxCLISessions[sessionID]; session != nil {
			session.Approving = false
		}
		h.entroxCLIMu.Unlock()
		response.ErrorFrom(c, err)
		return
	}

	h.entroxCLIMu.Lock()
	defer h.entroxCLIMu.Unlock()
	session = h.entroxCLISessions[sessionID]
	if session == nil {
		response.NotFound(c, "Auth session not found")
		return
	}
	if time.Now().After(session.ExpiresAt) {
		delete(h.entroxCLISessions, sessionID)
		response.Error(c, http.StatusGone, "Auth session expired")
		return
	}
	session.APIKey = key.Key
	session.APIKeyID = key.ID
	session.Created = created
	session.UserID = subject.UserID
	session.Approving = false

	response.Success(c, approveEntroxCLIResponse{
		Status:   "approved",
		APIKeyID: key.ID,
		Created:  created,
	})
}

func validateApproveEntroxCLIRequest(req approveEntroxCLIRequest) error {
	hasKeyID := req.APIKeyID != nil
	if hasKeyID && *req.APIKeyID <= 0 {
		return infraerrors.BadRequest("ENTROX_API_KEY_INVALID", "api key id must be positive")
	}
	if hasKeyID == req.CreateNew {
		return infraerrors.BadRequest(
			"ENTROX_API_KEY_CHOICE_REQUIRED",
			"choose an existing api key or create a new one",
		)
	}
	return nil
}

func (h *AuthHandler) resolveEntroxCLIAPIKey(ctx context.Context, userID int64, req approveEntroxCLIRequest, now time.Time) (*service.APIKey, bool, error) {
	if req.APIKeyID != nil {
		key, err := h.apiKeyService.GetByID(ctx, *req.APIKeyID)
		if err != nil {
			return nil, false, err
		}
		if key.UserID != userID {
			return nil, false, service.ErrInsufficientPerms
		}
		if key.Key == "" || !key.IsActive() || key.IsExpired() || key.IsQuotaExhausted() {
			return nil, false, infraerrors.BadRequest("ENTROX_API_KEY_UNAVAILABLE", "api key is not available")
		}
		return key, false, nil
	}

	key, err := h.apiKeyService.Create(ctx, userID, service.CreateAPIKeyRequest{
		Name: "entrox CLI " + now.Format("2006-01-02 15:04"),
	})
	if err != nil {
		return nil, false, err
	}
	return key, true, nil
}

func (h *AuthHandler) cleanupExpiredEntroxCLISessionsLocked(now time.Time) {
	for id, session := range h.entroxCLISessions {
		if now.After(session.ExpiresAt) {
			delete(h.entroxCLISessions, id)
		}
	}
}

func randomHex(bytesLen int) (string, error) {
	buf := make([]byte, bytesLen)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func authRequestOrigin(c *gin.Context) string {
	scheme := firstForwardedHeaderValue(c.GetHeader("X-Forwarded-Proto"))
	if scheme == "" {
		if c.Request.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}

	host := firstForwardedHeaderValue(c.GetHeader("X-Forwarded-Host"))
	if host == "" {
		host = c.Request.Host
	}
	if host == "" {
		host = "localhost"
	}

	return scheme + "://" + host
}

func firstForwardedHeaderValue(value string) string {
	if i := strings.IndexByte(value, ','); i >= 0 {
		value = value[:i]
	}
	return strings.TrimSpace(value)
}

const entroxCLILoginScript = `#!/bin/sh
set -eu

base="${1%/}"

fail() {
  printf '%s\n' "$*" >&2
  exit 1
}

json_field() {
  printf '%s' "$1" | sed -n "s/.*\"$2\":\"\([^\"]*\)\".*/\1/p" | head -n 1
}

command -v curl >/dev/null 2>&1 || fail "curl is required for entrox login"

start="$(curl -fsS -X POST "$base/api/v1/auth/entrox/start")" || fail "Failed to start entrox login"
session_id="$(json_field "$start" session_id)"
poll_token="$(json_field "$start" poll_token)"
authorize_url="$(json_field "$start" authorize_url)"

[ -n "$session_id" ] || fail "Missing entrox login session"
[ -n "$poll_token" ] || fail "Missing entrox poll token"
[ -n "$authorize_url" ] || fail "Missing entrox authorize URL"

if command -v python3 >/dev/null 2>&1; then
  python3 -c 'import sys, webbrowser; webbrowser.open(sys.argv[1])' "$authorize_url" >/dev/null 2>&1 || true
elif command -v open >/dev/null 2>&1; then
  open "$authorize_url" >/dev/null 2>&1 || true
elif command -v xdg-open >/dev/null 2>&1; then
  xdg-open "$authorize_url" >/dev/null 2>&1 || true
fi

printf 'Complete Sub2API login in your browser: %s\n' "$authorize_url" >&2

attempt=0
while [ "$attempt" -lt 300 ]; do
  body="$(curl -sS "$base/api/v1/auth/entrox/poll?session_id=$session_id&poll_token=$poll_token" || true)"
  token="$(json_field "$body" api_key)"
  if [ -n "$token" ]; then
    printf '%s' "$token"
    exit 0
  fi
  status="$(json_field "$body" status)"
  case "$status" in
    pending|"")
      sleep 2
      attempt=$((attempt + 1))
      ;;
    expired)
      fail "entrox login expired"
      ;;
    not_found|invalid_request)
      fail "entrox login session was not found"
      ;;
    *)
      fail "entrox login failed: $status"
      ;;
  esac
done

fail "Timed out waiting for browser authorization"
`
