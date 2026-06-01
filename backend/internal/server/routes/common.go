package routes

import (
	"context"
	_ "embed"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const sub2APIOpenCodeTokenEnv = "SUB2API_API_KEY"
const sub2APIOpenCodeProviderID = "entrox"
const sub2APIOpenCodeDynamicConfigPath = "/api/v1/entrox/opencode/config"

type entroxDownloadMirrorResolver interface {
	GetEntroxDownloadMirrorBaseURL(ctx context.Context) (string, error)
}

//go:embed entrox_install.sh
var entroxInstallScript string

//go:embed entrox_install.ps1
var entroxInstallPowerShellScript string

// RegisterCommonRoutes 注册通用路由（健康检查、状态等）
func RegisterCommonRoutes(r *gin.Engine, resolvers ...entroxDownloadMirrorResolver) {
	var mirrorResolver entroxDownloadMirrorResolver
	if len(resolvers) > 0 {
		mirrorResolver = resolvers[0]
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/install", func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache")
		c.Data(http.StatusOK, "text/x-shellscript; charset=utf-8", []byte(entroxInstallScript))
	})

	r.GET("/install.ps1", func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache")
		c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(entroxInstallPowerShellScript))
	})

	r.GET("/downloads/entrox-dev/:asset", func(c *gin.Context) {
		asset := c.Param("asset")
		if !isEntroxDevReleaseAsset(asset) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		baseURL, err := entroxDownloadMirrorBaseURL(c.Request.Context(), mirrorResolver)
		if err != nil {
			c.String(http.StatusInternalServerError, "Entrox download mirror is invalid")
			return
		}
		if baseURL == "" {
			c.String(http.StatusServiceUnavailable, "Entrox download mirror is not configured")
			return
		}
		c.Redirect(http.StatusTemporaryRedirect, baseURL+"/"+asset)
	})

	r.GET("/.well-known/opencode", func(c *gin.Context) {
		origin := requestOrigin(c)

		c.JSON(http.StatusOK, gin.H{
			"auth": gin.H{
				"command": []string{
					"sh",
					"-c",
					`curl -fsSL "$1/api/v1/auth/entrox/cli.sh" | sh -s -- "$1"`,
					"entrox-login",
					origin,
				},
				"env": sub2APIOpenCodeTokenEnv,
			},
			"config": gin.H{
				"$schema":     "https://opencode.ai/config.json",
				"model":       sub2APIOpenCodeProviderID + "/gpt-5.4",
				"small_model": sub2APIOpenCodeProviderID + "/gpt-5.4-mini",
			},
			"remote_config": gin.H{
				"url": origin + sub2APIOpenCodeDynamicConfigPath,
				"headers": gin.H{
					"Authorization": "Bearer {env:" + sub2APIOpenCodeTokenEnv + "}",
				},
			},
		})
	})

	// Claude Code 遥测日志（忽略，直接返回200）
	r.POST("/api/event_logging/batch", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Setup status endpoint (always returns needs_setup: false in normal mode)
	// This is used by the frontend to detect when the service has restarted after setup
	r.GET("/setup/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": gin.H{
				"needs_setup": false,
				"step":        "completed",
			},
		})
	})
}

func requestOrigin(c *gin.Context) string {
	scheme := firstHeaderValue(c.GetHeader("X-Forwarded-Proto"))
	if scheme == "" {
		if c.Request.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}

	host := firstHeaderValue(c.GetHeader("X-Forwarded-Host"))
	if host == "" {
		host = c.Request.Host
	}
	if host == "" {
		host = "localhost"
	}

	return scheme + "://" + host
}

func firstHeaderValue(value string) string {
	if i := strings.IndexByte(value, ','); i >= 0 {
		value = value[:i]
	}
	return strings.TrimSpace(value)
}

func entroxDownloadMirrorBaseURL(ctx context.Context, resolver entroxDownloadMirrorResolver) (string, error) {
	if resolver != nil {
		baseURL, err := resolver.GetEntroxDownloadMirrorBaseURL(ctx)
		if err != nil {
			return "", err
		}
		if baseURL != "" {
			return normalizeEntroxDownloadMirrorBaseURL(baseURL)
		}
	}
	return normalizeEntroxDownloadMirrorBaseURL(os.Getenv("ENTROX_DOWNLOAD_MIRROR_BASE_URL"))
}

func normalizeEntroxDownloadMirrorBaseURL(raw string) (string, error) {
	raw = strings.TrimRight(strings.TrimSpace(raw), "/")
	if raw == "" {
		return "", nil
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	if parsed.Host == "" || (parsed.Scheme != "https" && parsed.Scheme != "http") {
		return "", errors.New("download mirror must be an absolute http(s) URL")
	}
	if parsed.RawQuery != "" || parsed.Fragment != "" {
		return "", errors.New("download mirror must not include query or fragment")
	}
	return raw, nil
}

func isEntroxDevReleaseAsset(asset string) bool {
	switch asset {
	case "entrox-cli-macos-arm64.zip", "entrox-cli-linux-x64.zip", "entrox-cli-windows-x64.zip":
		return true
	default:
		return false
	}
}
