package handler

import (
	"net/http"
	"sort"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/claude"
	"github.com/Wei-Shaw/sub2api/internal/pkg/geminicli"
	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

const (
	entroxOpenCodeTokenEnv          = "SUB2API_API_KEY"
	entroxOpenCodeProviderID        = "entrox"
	entroxOpenCodeAnthropicProvider = "entrox-anthropic"
	entroxOpenCodeGeminiProvider    = "entrox-gemini"
	entroxOpenCodeProviderName      = "Entrox"
	entroxOpenCodeDefaultModel      = "gpt-5.4"
	entroxOpenCodeDefaultSmallModel = "gpt-5.4-mini"
	entroxOpenCodeOpenAINPM         = "@ai-sdk/openai-compatible"
	entroxOpenCodeAnthropicNPM      = "@ai-sdk/anthropic"
	entroxOpenCodeGeminiNPM         = "@ai-sdk/google"
)

// EntroxOpenCodeConfig returns the opencode provider config visible to the current API key.
func (h *GatewayHandler) EntroxOpenCodeConfig(c *gin.Context) {
	apiKey, ok := middleware2.GetAPIKeyFromContext(c)
	if !ok || apiKey == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
		return
	}
	if apiKey.Group == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "API key is not assigned to a group"})
		return
	}

	origin := entroxOpenCodeRequestOrigin(c)
	config := h.buildEntroxOpenCodeConfig(
		origin,
		h.gatewayService.GetAvailableModelsByPlatform(c.Request.Context(), &apiKey.Group.ID),
		apiKey.Group,
	)
	c.JSON(http.StatusOK, gin.H{"config": config})
}

func (h *GatewayHandler) buildEntroxOpenCodeConfig(origin string, modelsByPlatform map[string][]string, group *service.Group) gin.H {
	providers := gin.H{}
	apiURL := origin + "/v1"
	geminiURL := origin + "/v1beta"

	openAIModels := h.entroxOpenCodeModelsForPlatform(modelsByPlatform, group, service.PlatformOpenAI)
	if len(openAIModels) > 0 {
		providers[entroxOpenCodeProviderID] = entroxOpenCodeProvider(entroxOpenCodeOpenAINPM, apiURL, openAIModels)
	}

	anthropicModels := h.entroxOpenCodeModelsForPlatform(modelsByPlatform, group, service.PlatformAnthropic)
	if len(anthropicModels) > 0 {
		providers[entroxOpenCodeAnthropicProvider] = entroxOpenCodeProvider(entroxOpenCodeAnthropicNPM, apiURL, anthropicModels)
	}

	geminiModels := h.entroxOpenCodeModelsForPlatform(modelsByPlatform, group, service.PlatformGemini)
	if len(geminiModels) > 0 {
		providers[entroxOpenCodeGeminiProvider] = entroxOpenCodeProvider(entroxOpenCodeGeminiNPM, geminiURL, geminiModels)
	}

	return gin.H{
		"$schema":     "https://opencode.ai/config.json",
		"model":       entroxOpenCodeDefaultProviderModel(providers),
		"small_model": entroxOpenCodeSmallProviderModel(providers),
		"provider":    providers,
	}
}

func (h *GatewayHandler) entroxOpenCodeModelsForPlatform(modelsByPlatform map[string][]string, group *service.Group, platform string) gin.H {
	models, platformConfigured := modelsByPlatform[platform]
	if !platformConfigured {
		return gin.H{}
	}
	if group != nil && group.CustomModelsListEnabled() {
		models = filterModelsByCustomList(models, defaultModelIDsForPlatform(platform), group.ModelsListConfig.Models)
	}
	if len(models) == 0 {
		models = defaultModelIDsForPlatform(platform)
	}

	out := gin.H{}
	for _, modelID := range models {
		modelID = strings.TrimSpace(modelID)
		if modelID == "" || entroxOpenCodeHideModel(modelID) {
			continue
		}
		out[modelID] = entroxOpenCodeModel(modelID, platform)
	}
	return out
}

func entroxOpenCodeProvider(npm string, apiURL string, models gin.H) gin.H {
	return gin.H{
		"name": entroxOpenCodeProviderName,
		"npm":  npm,
		"api":  apiURL,
		"options": gin.H{
			"baseURL": apiURL,
			"apiKey":  "{env:" + entroxOpenCodeTokenEnv + "}",
		},
		"models": models,
	}
}

func entroxOpenCodeModel(modelID string, platform string) gin.H {
	return gin.H{
		"name":       entroxOpenCodeDisplayName(modelID, platform),
		"tool_call":  true,
		"reasoning":  entroxOpenCodeReasoning(modelID, platform),
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

func entroxOpenCodeDisplayName(modelID string, platform string) string {
	switch platform {
	case service.PlatformOpenAI:
		for _, model := range openai.DefaultModels {
			if model.ID == modelID {
				return model.DisplayName
			}
		}
	case service.PlatformAnthropic:
		for _, model := range claude.DefaultModels {
			if model.ID == modelID {
				return model.DisplayName
			}
		}
	case service.PlatformGemini:
		for _, model := range geminicli.DefaultModels {
			if model.ID == modelID {
				return model.DisplayName
			}
		}
	}
	return modelID
}

func entroxOpenCodeReasoning(modelID string, platform string) bool {
	if platform == service.PlatformOpenAI {
		return strings.HasPrefix(modelID, "gpt-5")
	}
	return strings.Contains(modelID, "thinking") || strings.Contains(modelID, "-pro")
}

func entroxOpenCodeHideModel(modelID string) bool {
	lower := strings.ToLower(modelID)
	return strings.Contains(lower, "image") || strings.Contains(lower, "embedding")
}

func entroxOpenCodeDefaultProviderModel(providers gin.H) string {
	if provider, ok := providers[entroxOpenCodeProviderID].(gin.H); ok {
		if models, ok := provider["models"].(gin.H); ok {
			if _, ok := models[entroxOpenCodeDefaultModel]; ok {
				return entroxOpenCodeProviderID + "/" + entroxOpenCodeDefaultModel
			}
			if modelID := firstSortedModelID(models); modelID != "" {
				return entroxOpenCodeProviderID + "/" + modelID
			}
		}
	}
	for _, providerID := range []string{entroxOpenCodeAnthropicProvider, entroxOpenCodeGeminiProvider} {
		provider, _ := providers[providerID].(gin.H)
		models, _ := provider["models"].(gin.H)
		if modelID := firstSortedModelID(models); modelID != "" {
			return providerID + "/" + modelID
		}
	}
	return entroxOpenCodeProviderID + "/" + entroxOpenCodeDefaultModel
}

func entroxOpenCodeSmallProviderModel(providers gin.H) string {
	if provider, ok := providers[entroxOpenCodeProviderID].(gin.H); ok {
		if models, ok := provider["models"].(gin.H); ok {
			if _, ok := models[entroxOpenCodeDefaultSmallModel]; ok {
				return entroxOpenCodeProviderID + "/" + entroxOpenCodeDefaultSmallModel
			}
		}
	}
	return entroxOpenCodeDefaultProviderModel(providers)
}

func firstSortedModelID(models gin.H) string {
	if len(models) == 0 {
		return ""
	}
	keys := make([]string, 0, len(models))
	for key := range models {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys[0]
}

func entroxOpenCodeRequestOrigin(c *gin.Context) string {
	scheme := entroxOpenCodeFirstHeaderValue(c.GetHeader("X-Forwarded-Proto"))
	if scheme == "" {
		if c.Request.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}

	host := entroxOpenCodeFirstHeaderValue(c.GetHeader("X-Forwarded-Host"))
	if host == "" {
		host = c.Request.Host
	}
	if host == "" {
		host = "localhost"
	}
	return scheme + "://" + host
}

func entroxOpenCodeFirstHeaderValue(value string) string {
	if i := strings.IndexByte(value, ','); i >= 0 {
		value = value[:i]
	}
	return strings.TrimSpace(value)
}
