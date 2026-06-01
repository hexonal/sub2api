package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestRegisterCommonRoutesEntroxInstallRedirect(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	RegisterCommonRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/install", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}
	contentType := rec.Header().Get("Content-Type")
	if contentType != "text/x-shellscript; charset=utf-8" {
		t.Fatalf("expected shell script content type, got %q", contentType)
	}
	body := rec.Body.String()
	if len(body) == 0 || body[:19] != "#!/usr/bin/env bash" {
		t.Fatalf("expected bash installer, got %q", body)
	}
	if !strings.Contains(body, "APP=entrox") {
		t.Fatalf("expected Entrox installer, got %q", body)
	}
	if !strings.Contains(body, "ENTROX_DOWNLOAD_BASE_URL") {
		t.Fatalf("expected mirror-aware installer, got %q", body)
	}
}

func TestRegisterCommonRoutesEntroxInstallPowerShell(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	RegisterCommonRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/install.ps1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}
	contentType := rec.Header().Get("Content-Type")
	if contentType != "text/plain; charset=utf-8" {
		t.Fatalf("expected powershell script content type, got %q", contentType)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "$App = \"entrox\"") {
		t.Fatalf("expected Entrox PowerShell installer, got %q", body)
	}
	if !strings.Contains(body, "ENTROX_DOWNLOAD_BASE_URL") {
		t.Fatalf("expected mirror-aware PowerShell installer, got %q", body)
	}
}

func TestRegisterCommonRoutesEntroxDevDownloadRedirectUsesMirror(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("ENTROX_DOWNLOAD_MIRROR_BASE_URL", "https://oss.example.test/entrox-dev/")

	router := gin.New()
	RegisterCommonRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/downloads/entrox-dev/entrox-cli-windows-x64.zip", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusTemporaryRedirect {
		t.Fatalf("expected redirect status, got %d: %s", rec.Code, rec.Body.String())
	}
	expected := "https://oss.example.test/entrox-dev/entrox-cli-windows-x64.zip"
	if rec.Header().Get("Location") != expected {
		t.Fatalf("expected redirect to %q, got %q", expected, rec.Header().Get("Location"))
	}
}

func TestRegisterCommonRoutesEntroxDevDownloadRejectsUnknownAsset(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	RegisterCommonRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/downloads/entrox-dev/other.zip", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestIsEntroxDevReleaseAsset(t *testing.T) {
	for _, asset := range []string{
		"entrox-cli-macos-arm64.zip",
		"entrox-cli-linux-x64.zip",
		"entrox-cli-windows-x64.zip",
	} {
		if !isEntroxDevReleaseAsset(asset) {
			t.Fatalf("expected %q to be allowed", asset)
		}
	}

	for _, asset := range []string{"other.zip", "../entrox-cli-linux-x64.zip", "entrox-cli-linux-x64.tar.gz"} {
		if isEntroxDevReleaseAsset(asset) {
			t.Fatalf("expected %q to be rejected", asset)
		}
	}
}
