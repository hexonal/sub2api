package routes

import (
	_ "embed"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const sub2APIOpenCodeTokenEnv = "SUB2API_API_KEY"
const sub2APIOpenCodeProviderID = "entrox"
const sub2APIOpenCodeDynamicConfigPath = "/api/v1/entrox/opencode/config"
const entroxDevReleaseFallbackBaseURL = "https://github.com/hexonal/entrox/releases/download/entrox-dev"

//go:embed entrox_install.sh
var entroxInstallScript string

//go:embed entrox_install.ps1
var entroxInstallPowerShellScript string

// RegisterCommonRoutes 注册通用路由（健康检查、状态等）
func RegisterCommonRoutes(r *gin.Engine) {
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

		baseURL := strings.TrimRight(os.Getenv("ENTROX_DOWNLOAD_MIRROR_BASE_URL"), "/")
		if baseURL != "" {
			c.Redirect(http.StatusTemporaryRedirect, baseURL+"/"+asset)
			return
		}

		proxyEntroxDevReleaseAsset(c, entroxDevReleaseFallbackBaseURL+"/"+asset)
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

func proxyEntroxDevReleaseAsset(c *gin.Context, url string) {
	req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, url, nil)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	client := &http.Client{Timeout: 10 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for _, header := range []string{"Content-Type", "Content-Length", "ETag", "Last-Modified"} {
		if value := resp.Header.Get(header); value != "" {
			c.Header(header, value)
		}
	}
	c.Status(resp.StatusCode)
	_, _ = io.Copy(c.Writer, resp.Body)
}

func isEntroxDevReleaseAsset(asset string) bool {
	switch asset {
	case "entrox-cli-macos-arm64.zip", "entrox-cli-linux-x64.zip", "entrox-cli-windows-x64.zip":
		return true
	default:
		return false
	}
}
