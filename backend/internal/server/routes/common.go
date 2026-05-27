package routes

import (
	"net/http"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"

	"github.com/gin-gonic/gin"
)

const sub2APIOpenCodeTokenEnv = "SUB2API_API_KEY"
const sub2APIOpenCodeProviderID = "entrox"

// RegisterCommonRoutes 注册通用路由（健康检查、状态等）
func RegisterCommonRoutes(r *gin.Engine) {
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
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
				"provider": gin.H{
					sub2APIOpenCodeProviderID: gin.H{
						"name": "Entrox",
						"npm":  "@ai-sdk/openai-compatible",
						"api":  origin + "/v1",
						"options": gin.H{
							"baseURL": origin + "/v1",
							"apiKey":  "{env:" + sub2APIOpenCodeTokenEnv + "}",
						},
						"models": sub2APIOpenCodeModels(),
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

func sub2APIOpenCodeModels() gin.H {
	models := gin.H{}
	for _, model := range openai.DefaultModels {
		if strings.Contains(model.ID, "image") {
			continue
		}
		models[model.ID] = gin.H{
			"name":       model.DisplayName,
			"tool_call":  true,
			"reasoning":  strings.HasPrefix(model.ID, "gpt-5"),
			"attachment": true,
			"modalities": gin.H{
				"input":  []string{"text", "image"},
				"output": []string{"text"},
			},
			"limit": gin.H{
				"context": 200000,
				"output":  32768,
			},
		}
	}
	return models
}

func firstHeaderValue(value string) string {
	if i := strings.IndexByte(value, ','); i >= 0 {
		value = value[:i]
	}
	return strings.TrimSpace(value)
}
