package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterCommonRoutesWellKnownOpenCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	RegisterCommonRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/.well-known/opencode", nil)
	req.Host = "sub2api.example.test"
	req.Header.Set("X-Forwarded-Proto", "https")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	auth, ok := body["auth"].(map[string]any)
	if !ok {
		t.Fatalf("expected auth object, got %#v", body["auth"])
	}
	if auth["env"] != sub2APIOpenCodeTokenEnv {
		t.Fatalf("expected auth env %q, got %#v", sub2APIOpenCodeTokenEnv, auth["env"])
	}
	command, ok := auth["command"].([]any)
	if !ok || len(command) < 5 {
		t.Fatalf("expected auth command array, got %#v", auth["command"])
	}
	if command[2] != `curl -fsSL "$1/api/v1/auth/entrox/cli.sh" | sh -s -- "$1"` {
		t.Fatalf("expected entrox cli script command, got %#v", command[2])
	}
	if command[len(command)-1] != "https://sub2api.example.test" {
		t.Fatalf("expected command origin, got %#v", command[len(command)-1])
	}

	config, ok := body["config"].(map[string]any)
	if !ok {
		t.Fatalf("expected config object, got %#v", body["config"])
	}
	provider := config["provider"].(map[string]any)
	openai := provider["openai"].(map[string]any)
	options := openai["options"].(map[string]any)
	if options["baseURL"] != "https://sub2api.example.test/v1" {
		t.Fatalf("expected Sub2API baseURL, got %#v", options["baseURL"])
	}
	if options["apiKey"] != "{env:"+sub2APIOpenCodeTokenEnv+"}" {
		t.Fatalf("expected env apiKey, got %#v", options["apiKey"])
	}
}
