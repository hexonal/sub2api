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
	if config["model"] != sub2APIOpenCodeProviderID+"/gpt-5.4" {
		t.Fatalf("expected Entrox default model, got %#v", config["model"])
	}
	if config["small_model"] != sub2APIOpenCodeProviderID+"/gpt-5.4-mini" {
		t.Fatalf("expected Entrox small model, got %#v", config["small_model"])
	}
	if _, ok := config["provider"]; ok {
		t.Fatalf("expected well-known provider config to be loaded dynamically")
	}

	remoteConfig, ok := body["remote_config"].(map[string]any)
	if !ok {
		t.Fatalf("expected remote_config object, got %#v", body["remote_config"])
	}
	if remoteConfig["url"] != "https://sub2api.example.test/api/v1/entrox/opencode/config" {
		t.Fatalf("expected dynamic config URL, got %#v", remoteConfig["url"])
	}
	headers, ok := remoteConfig["headers"].(map[string]any)
	if !ok {
		t.Fatalf("expected remote_config headers object, got %#v", remoteConfig["headers"])
	}
	if headers["Authorization"] != "Bearer {env:"+sub2APIOpenCodeTokenEnv+"}" {
		t.Fatalf("expected auth header env substitution, got %#v", headers["Authorization"])
	}
}
