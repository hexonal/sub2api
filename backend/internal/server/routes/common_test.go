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
		t.Fatalf("expected Sub2API default model, got %#v", config["model"])
	}
	if config["small_model"] != sub2APIOpenCodeProviderID+"/gpt-5.4-mini" {
		t.Fatalf("expected Sub2API small model, got %#v", config["small_model"])
	}
	provider := config["provider"].(map[string]any)
	if _, ok := provider["openai"]; ok {
		t.Fatalf("expected well-known config not to override built-in openai provider")
	}
	sub2api := provider[sub2APIOpenCodeProviderID].(map[string]any)
	if sub2api["api"] != "https://sub2api.example.test/v1" {
		t.Fatalf("expected Sub2API api URL, got %#v", sub2api["api"])
	}
	if sub2api["npm"] != "@ai-sdk/openai-compatible" {
		t.Fatalf("expected openai-compatible provider, got %#v", sub2api["npm"])
	}
	options := sub2api["options"].(map[string]any)
	if options["baseURL"] != "https://sub2api.example.test/v1" {
		t.Fatalf("expected Sub2API baseURL, got %#v", options["baseURL"])
	}
	if options["apiKey"] != "{env:"+sub2APIOpenCodeTokenEnv+"}" {
		t.Fatalf("expected env apiKey, got %#v", options["apiKey"])
	}
	models := sub2api["models"].(map[string]any)
	if _, ok := models["gpt-5.4"]; !ok {
		t.Fatalf("expected gpt-5.4 model in Sub2API provider")
	}
	if _, ok := models["gpt-image-1"]; ok {
		t.Fatalf("expected image-only models to be hidden from entrox model picker")
	}
}
