package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestStartEntroxCLIAuthCreatesPollableSession(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewAuthHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	r := gin.New()
	r.POST("/start", h.StartEntroxCLIAuth)
	r.GET("/poll", h.PollEntroxCLIAuth)

	startReq := httptest.NewRequest(http.MethodPost, "/start", nil)
	startReq.Host = "sub2api.example.test"
	startReq.Header.Set("X-Forwarded-Proto", "https")
	startResp := httptest.NewRecorder()

	r.ServeHTTP(startResp, startReq)

	require.Equal(t, http.StatusOK, startResp.Code)
	var startBody struct {
		SessionID    string `json:"session_id"`
		PollToken    string `json:"poll_token"`
		AuthorizeURL string `json:"authorize_url"`
		ExpiresIn    int    `json:"expires_in"`
	}
	require.NoError(t, json.Unmarshal(startResp.Body.Bytes(), &startBody))
	require.NotEmpty(t, startBody.SessionID)
	require.NotEmpty(t, startBody.PollToken)
	require.Equal(t, "https://sub2api.example.test/entrox/connect?session_id="+startBody.SessionID, startBody.AuthorizeURL)
	require.Equal(t, int(entroxCLISessionTTL.Seconds()), startBody.ExpiresIn)

	pollReq := httptest.NewRequest(http.MethodGet, "/poll?session_id="+startBody.SessionID+"&poll_token="+startBody.PollToken, nil)
	pollResp := httptest.NewRecorder()

	r.ServeHTTP(pollResp, pollReq)

	require.Equal(t, http.StatusAccepted, pollResp.Code)
	require.JSONEq(t, `{"status":"pending"}`, pollResp.Body.String())
}

func TestEntroxCLIScript(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewAuthHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	r := gin.New()
	r.GET("/cli.sh", h.EntroxCLIScript)

	req := httptest.NewRequest(http.MethodGet, "/cli.sh", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	require.Contains(t, resp.Header().Get("Content-Type"), "text/x-shellscript")
	require.Contains(t, resp.Body.String(), "/api/v1/auth/entrox/start")
	require.Contains(t, resp.Body.String(), "/api/v1/auth/entrox/poll")
}

func TestApproveEntroxCLIAuthRequiresExplicitAPIKeyChoice(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := newTestEntroxCLIAuthHandler(7, &testEntroxAPIKeyRepo{})
	seedTestEntroxCLISession(h, "session-choice", "poll-choice")
	r := newTestEntroxCLIRouter(h, 7)

	resp := postTestEntroxCLIApprove(t, r, `{"session_id":"session-choice"}`)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	require.Contains(t, resp.Body.String(), "api key")

	pollReq := httptest.NewRequest(http.MethodGet, "/poll?session_id=session-choice&poll_token=poll-choice", nil)
	pollResp := httptest.NewRecorder()
	r.ServeHTTP(pollResp, pollReq)

	require.Equal(t, http.StatusAccepted, pollResp.Code)
	require.JSONEq(t, `{"status":"pending"}`, pollResp.Body.String())
}

func TestApproveEntroxCLIAuthUsesSelectedExistingAPIKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	groupID := int64(3)
	repo := &testEntroxAPIKeyRepo{}
	repo.seed(&service.APIKey{
		ID:      11,
		UserID:  7,
		Key:     "sk-existing",
		Name:    "Existing key",
		GroupID: &groupID,
		Status:  service.StatusActive,
	})
	h := newTestEntroxCLIAuthHandler(7, repo)
	seedTestEntroxCLISession(h, "session-existing", "poll-existing")
	r := newTestEntroxCLIRouter(h, 7)

	resp := postTestEntroxCLIApprove(t, r, `{"session_id":"session-existing","api_key_id":11}`)

	require.Equal(t, http.StatusOK, resp.Code)
	require.Len(t, repo.created, 0)

	var body struct {
		Data struct {
			Status   string `json:"status"`
			APIKeyID int64  `json:"api_key_id"`
			Created  bool   `json:"created"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	require.Equal(t, "approved", body.Data.Status)
	require.Equal(t, int64(11), body.Data.APIKeyID)
	require.False(t, body.Data.Created)

	pollReq := httptest.NewRequest(http.MethodGet, "/poll?session_id=session-existing&poll_token=poll-existing", nil)
	pollResp := httptest.NewRecorder()
	r.ServeHTTP(pollResp, pollReq)

	require.Equal(t, http.StatusOK, pollResp.Code)
	require.JSONEq(t, `{"status":"approved","api_key":"sk-existing"}`, pollResp.Body.String())
}

func TestApproveEntroxCLIAuthRejectsUnassignedAPIKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &testEntroxAPIKeyRepo{}
	repo.seed(&service.APIKey{
		ID:     11,
		UserID: 7,
		Key:    "sk-existing",
		Name:   "Existing key",
		Status: service.StatusActive,
	})
	h := newTestEntroxCLIAuthHandler(7, repo)
	seedTestEntroxCLISession(h, "session-unassigned", "poll-unassigned")
	r := newTestEntroxCLIRouter(h, 7)

	resp := postTestEntroxCLIApprove(t, r, `{"session_id":"session-unassigned","api_key_id":11}`)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	require.Contains(t, resp.Body.String(), "assigned to a group")

	pollReq := httptest.NewRequest(http.MethodGet, "/poll?session_id=session-unassigned&poll_token=poll-unassigned", nil)
	pollResp := httptest.NewRecorder()
	r.ServeHTTP(pollResp, pollReq)

	require.Equal(t, http.StatusAccepted, pollResp.Code)
	require.JSONEq(t, `{"status":"pending"}`, pollResp.Body.String())
}

func TestApproveEntroxCLIAuthRejectsNonEntroxAPIKeyGroup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	groupID := int64(4)
	repo := &testEntroxAPIKeyRepo{}
	repo.seed(&service.APIKey{
		ID:      11,
		UserID:  7,
		Key:     "sk-existing",
		Name:    "Existing key",
		GroupID: &groupID,
		Status:  service.StatusActive,
	})
	h := newTestEntroxCLIAuthHandler(7, repo)
	seedTestEntroxCLISession(h, "session-wrong-platform", "poll-wrong-platform")
	r := newTestEntroxCLIRouter(h, 7)

	resp := postTestEntroxCLIApprove(t, r, `{"session_id":"session-wrong-platform","api_key_id":11}`)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	require.Contains(t, resp.Body.String(), "entrox group")
}

func TestApproveEntroxCLIAuthCreatesAPIKeyOnlyWhenRequested(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &testEntroxAPIKeyRepo{}
	h := newTestEntroxCLIAuthHandler(7, repo)
	seedTestEntroxCLISession(h, "session-new", "poll-new")
	r := newTestEntroxCLIRouter(h, 7)

	resp := postTestEntroxCLIApprove(t, r, `{"session_id":"session-new","create_new":true,"group_id":3}`)

	require.Equal(t, http.StatusOK, resp.Code)
	require.Len(t, repo.created, 1)
	require.Equal(t, int64(7), repo.created[0].UserID)
	require.NotNil(t, repo.created[0].GroupID)
	require.Equal(t, int64(3), *repo.created[0].GroupID)
	require.Equal(t, service.StatusActive, repo.created[0].Status)
	require.NotEmpty(t, repo.created[0].Key)

	var body struct {
		Data struct {
			Status   string `json:"status"`
			APIKeyID int64  `json:"api_key_id"`
			Created  bool   `json:"created"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	require.Equal(t, "approved", body.Data.Status)
	require.Equal(t, repo.created[0].ID, body.Data.APIKeyID)
	require.True(t, body.Data.Created)
}

func TestApproveEntroxCLIAuthRejectsNewAPIKeyForNonEntroxGroup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &testEntroxAPIKeyRepo{}
	h := newTestEntroxCLIAuthHandler(7, repo)
	seedTestEntroxCLISession(h, "session-new-wrong-platform", "poll-new-wrong-platform")
	r := newTestEntroxCLIRouter(h, 7)

	resp := postTestEntroxCLIApprove(t, r, `{"session_id":"session-new-wrong-platform","create_new":true,"group_id":4}`)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	require.Contains(t, resp.Body.String(), "entrox group")
	require.Len(t, repo.created, 0)
}

func newTestEntroxCLIAuthHandler(userID int64, repo *testEntroxAPIKeyRepo) *AuthHandler {
	cfg := &config.Config{}
	cfg.Default.APIKeyPrefix = "sk-"
	userRepo := &testEntroxUserRepo{
		user: &service.User{ID: userID, Email: "user@example.test", Status: service.StatusActive},
	}
	groupRepo := &testEntroxGroupRepo{
		byID: map[int64]*service.Group{
			3: {
				ID:               3,
				Name:             "Entrox Pro",
				Platform:         service.PlatformEntrox,
				Status:           service.StatusActive,
				SubscriptionType: service.SubscriptionTypeStandard,
			},
			4: {
				ID:               4,
				Name:             "OpenAI Pro",
				Platform:         service.PlatformOpenAI,
				Status:           service.StatusActive,
				SubscriptionType: service.SubscriptionTypeStandard,
			},
		},
	}
	apiKeySvc := service.NewAPIKeyService(repo, userRepo, groupRepo, nil, nil, nil, cfg)
	return NewAuthHandler(nil, nil, nil, nil, nil, nil, nil, apiKeySvc, nil)
}

func seedTestEntroxCLISession(h *AuthHandler, sessionID string, pollToken string) {
	now := time.Now()
	h.entroxCLISessions[sessionID] = &entroxCLISession{
		ID:        sessionID,
		PollToken: pollToken,
		CreatedAt: now,
		ExpiresAt: now.Add(entroxCLISessionTTL),
	}
}

func newTestEntroxCLIRouter(h *AuthHandler, userID int64) *gin.Engine {
	r := gin.New()
	r.POST("/approve", func(c *gin.Context) {
		c.Set(string(middleware2.ContextKeyUser), middleware2.AuthSubject{UserID: userID})
		h.ApproveEntroxCLIAuth(c)
	})
	r.GET("/poll", h.PollEntroxCLIAuth)
	return r
}

func postTestEntroxCLIApprove(t *testing.T, r *gin.Engine, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/approve", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	return resp
}

type testEntroxUserRepo struct {
	service.UserRepository
	user *service.User
}

func (r *testEntroxUserRepo) GetByID(ctx context.Context, id int64) (*service.User, error) {
	if r.user == nil || r.user.ID != id {
		return nil, service.ErrUserNotFound
	}
	clone := *r.user
	return &clone, nil
}

type testEntroxAPIKeyRepo struct {
	service.APIKeyRepository
	nextID  int64
	byID    map[int64]*service.APIKey
	created []*service.APIKey
}

func (r *testEntroxAPIKeyRepo) ensure() {
	if r.nextID == 0 {
		r.nextID = 100
	}
	if r.byID == nil {
		r.byID = make(map[int64]*service.APIKey)
	}
}

func (r *testEntroxAPIKeyRepo) seed(key *service.APIKey) {
	r.ensure()
	clone := *key
	r.byID[clone.ID] = &clone
}

func (r *testEntroxAPIKeyRepo) Create(ctx context.Context, key *service.APIKey) error {
	r.ensure()
	if key.ID == 0 {
		key.ID = r.nextID
		r.nextID++
	}
	now := time.Now()
	if key.CreatedAt.IsZero() {
		key.CreatedAt = now
	}
	if key.UpdatedAt.IsZero() {
		key.UpdatedAt = now
	}
	clone := *key
	r.byID[clone.ID] = &clone
	r.created = append(r.created, &clone)
	return nil
}

func (r *testEntroxAPIKeyRepo) GetByID(ctx context.Context, id int64) (*service.APIKey, error) {
	r.ensure()
	key := r.byID[id]
	if key == nil {
		return nil, service.ErrAPIKeyNotFound
	}
	clone := *key
	return &clone, nil
}

func (r *testEntroxAPIKeyRepo) ListByUserID(ctx context.Context, userID int64, params pagination.PaginationParams, filters service.APIKeyListFilters) ([]service.APIKey, *pagination.PaginationResult, error) {
	return nil, nil, nil
}

type testEntroxGroupRepo struct {
	service.GroupRepository
	byID map[int64]*service.Group
}

func (r *testEntroxGroupRepo) GetByID(ctx context.Context, id int64) (*service.Group, error) {
	group := r.byID[id]
	if group == nil {
		return nil, service.ErrGroupNotFound
	}
	clone := *group
	return &clone, nil
}
