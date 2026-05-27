package routes

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const sub2APIOpenCodeTokenEnv = "SUB2API_API_KEY"

// RegisterCommonRoutes 注册通用路由（健康检查、状态等）
func RegisterCommonRoutes(r *gin.Engine) {
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/.well-known/opencode", func(c *gin.Context) {
		origin := requestOrigin(c)
		authScript := strings.Join([]string{
			`login_url="$1/keys"`,
			`if command -v python3 >/dev/null 2>&1; then`,
			`  python3 -c 'import sys, webbrowser; webbrowser.open(sys.argv[1])' "$login_url" >/dev/null 2>&1 || true`,
			`elif command -v open >/dev/null 2>&1; then`,
			`  open "$login_url" >/dev/null 2>&1 || true`,
			`elif command -v xdg-open >/dev/null 2>&1; then`,
			`  xdg-open "$login_url" >/dev/null 2>&1 || true`,
			`fi`,
			`printf 'Paste Sub2API API key from %s: ' "$login_url" >/dev/tty`,
			`IFS= read -r token </dev/tty`,
			`printf '%s' "$token"`,
		}, "\n")

		c.JSON(http.StatusOK, gin.H{
			"auth": gin.H{
				"command": []string{"sh", "-c", authScript, "sub2api-entrox-login", origin},
				"env":     sub2APIOpenCodeTokenEnv,
			},
			"config": gin.H{
				"$schema":     "https://opencode.ai/config.json",
				"model":       "openai/gpt-5.4",
				"small_model": "openai/gpt-5.4-mini",
				"provider": gin.H{
					"openai": gin.H{
						"name": "Sub2API OpenAI",
						"options": gin.H{
							"baseURL": origin + "/v1",
							"apiKey":  "{env:" + sub2APIOpenCodeTokenEnv + "}",
						},
					},
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
