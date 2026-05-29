package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
