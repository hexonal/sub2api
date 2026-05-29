package routes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/handler"
	"github.com/Wei-Shaw/sub2api/internal/pkg/ctxkey"
	servermiddleware "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newGatewayRoutesTestRouter(groupPlatform string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	RegisterGatewayRoutes(
		router,
		&handler.Handlers{
			Gateway:       &handler.GatewayHandler{},
			OpenAIGateway: &handler.OpenAIGatewayHandler{},
		},
		servermiddleware.APIKeyAuthMiddleware(func(c *gin.Context) {
			groupID := int64(1)
			group := &service.Group{
				ID:       groupID,
				Platform: groupPlatform,
				Status:   service.StatusActive,
				Hydrated: true,
			}
			c.Set(string(servermiddleware.ContextKeyAPIKey), &service.APIKey{
				GroupID: &groupID,
				Group:   group,
			})
			ctx := context.WithValue(c.Request.Context(), ctxkey.Group, group)
			c.Request = c.Request.WithContext(ctx)
			c.Next()
		}),
		nil,
		nil,
		nil,
		nil,
		&config.Config{},
	)

	return router
}

func TestGatewayRoutesOpenAIResponsesCompactPathIsRegistered(t *testing.T) {
	router := newGatewayRoutesTestRouter(service.PlatformOpenAI)

	for _, path := range []string{
		"/v1/responses/compact",
		"/responses/compact",
		"/backend-api/codex/responses",
		"/backend-api/codex/responses/compact",
	} {
		req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(`{"model":"gpt-5"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		require.NotEqual(t, http.StatusNotFound, w.Code, "path=%s should hit OpenAI responses handler", path)
	}
}

func TestGatewayRoutesOpenAIImagesPathsAreRegistered(t *testing.T) {
	router := newGatewayRoutesTestRouter(service.PlatformOpenAI)

	for _, path := range []string{
		"/v1/images/generations",
		"/v1/images/edits",
		"/images/generations",
		"/images/edits",
	} {
		req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(`{"model":"gpt-image-2","prompt":"draw a cat"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		require.NotEqual(t, http.StatusNotFound, w.Code, "path=%s should hit OpenAI images handler", path)
	}
}

func TestGatewayRoutesEntroxResolvesPlatformByEndpoint(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "openai responses", path: "/v1/responses", want: service.PlatformOpenAI},
		{name: "openai responses alias", path: "/responses", want: service.PlatformOpenAI},
		{name: "codex responses alias", path: "/backend-api/codex/responses", want: service.PlatformOpenAI},
		{name: "openai chat completions alias", path: "/chat/completions", want: service.PlatformOpenAI},
		{name: "openai embeddings alias", path: "/embeddings", want: service.PlatformOpenAI},
		{name: "openai images alias", path: "/images/generations", want: service.PlatformOpenAI},
		{name: "gemini native", path: "/v1beta/models/gemini-2.5-flash:generateContent", want: service.PlatformGemini},
		{name: "anthropic messages", path: "/v1/messages", want: service.PlatformAnthropic},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			req := httptest.NewRequest(http.MethodPost, tt.path, strings.NewReader(`{"model":"test"}`))
			groupID := int64(1)
			group := &service.Group{
				ID:       groupID,
				Platform: service.PlatformEntrox,
				Status:   service.StatusActive,
				Hydrated: true,
			}
			req = req.WithContext(context.WithValue(req.Context(), ctxkey.Group, group))
			c.Request = req
			c.Set(string(servermiddleware.ContextKeyAPIKey), &service.APIKey{
				GroupID: &groupID,
				Group:   group,
			})

			got := getGroupPlatform(c)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.want, c.Request.Context().Value(ctxkey.RequestPlatform))
		})
	}
}

func TestResolveGroupPlatformMiddlewareSetsEntroxGeminiPlatform(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	req := httptest.NewRequest(http.MethodPost, "/v1beta/models/gemini-2.5-flash:generateContent", strings.NewReader(`{}`))
	groupID := int64(1)
	group := &service.Group{
		ID:       groupID,
		Platform: service.PlatformEntrox,
		Status:   service.StatusActive,
		Hydrated: true,
	}
	req = req.WithContext(context.WithValue(req.Context(), ctxkey.Group, group))
	c.Request = req
	c.Set(string(servermiddleware.ContextKeyAPIKey), &service.APIKey{
		GroupID: &groupID,
		Group:   group,
	})

	resolveGroupPlatformMiddleware()(c)

	require.Equal(t, service.PlatformGemini, c.Request.Context().Value(ctxkey.RequestPlatform))
	require.Equal(t, service.PlatformGemini, c.Request.Context().Value(ctxkey.Platform))
}
