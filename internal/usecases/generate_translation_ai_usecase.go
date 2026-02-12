package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/config"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/errors"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// GenerateTranslationAIUseCase defines the interface for AI translation operations
type GenerateTranslationAIUseCase interface {
	TranslateCurriculum(ctx context.Context, req *dto.GenerateTranslationAIRequest) (*dto.GenerateTranslationAIResponse, error)
}

// generateTranslationAIUseCase implements GenerateTranslationAIUseCase interface
type generateTranslationAIUseCase struct {
	openaiClient *openai.Client
}

// NewGenerateIntroAIUseCase creates a new instance of GenerateTranslationAIUseCase
func NewGenerateTranslationAIUseCase(apiKey string) (GenerateTranslationAIUseCase, error) {
	if apiKey == "" {
		return nil, errors.NewAppError("OPENAI_API_KEY environment variable is required")
	}

	client := openai.NewClient(option.WithAPIKey(apiKey))

	return &generateTranslationAIUseCase{
		openaiClient: &client,
	}, nil
}

// TranslateCurriculum translates the curriculum data to the target language
func (uc *generateTranslationAIUseCase) TranslateCurriculum(ctx context.Context, req *dto.GenerateTranslationAIRequest) (*dto.GenerateTranslationAIResponse, error) {
	// Convert curriculum data to JSON string for processing
	curriculumJSON, err := json.Marshal(req.CurriculumData)
	if err != nil {
		return nil, errors.WrapError(err, "failed to marshal curriculum data")
	}

	// Get target language name
	targetLanguage := uc.getLanguageName(req.TargetLanguage)

	// Prepare the prompt for translation
	prompt := fmt.Sprintf(`
Translate the following curriculum JSON to %s. 
You must translate ALL text fields (strings) while keeping the JSON structure exactly the same.
Do NOT translate field names/keys, only translate the values that are strings.
Keep all dates, IDs, and technical fields unchanged.
Return ONLY the translated JSON, no additional text or explanations.

Curriculum JSON:
%s
`, targetLanguage, string(curriculumJSON))

	// Create chat completion request
	chatReq := openai.ChatCompletionNewParams{
		Model: "gpt-4o-mini",
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(`You are a professional translator specialized in translating curriculum vitae and professional documents. 
Your task is to translate curriculum data while maintaining the exact JSON structure and only translating text content.
Rules:
1. Translate ONLY string values, never field names/keys
2. Keep all dates, IDs, numbers, and technical fields unchanged
3. Maintain the exact JSON structure
4. Return ONLY the translated JSON, no explanations
5. Preserve formatting and special characters where appropriate`),
			openai.UserMessage(prompt),
		},
		MaxTokens:   openai.Int(int64(config.ParseIntEnv("OPENAI_MAX_TOKENS", 4000))),
		Temperature: openai.Float(config.ParseFloatEnv("OPENAI_TEMPERATURE", 0.3)),
		TopP:        openai.Float(config.ParseFloatEnv("OPENAI_TOP_P", 1.0)),
	}

	// Call OpenAI API
	resp, err := uc.openaiClient.Chat.Completions.New(ctx, chatReq)
	if err != nil {
		return nil, errors.WrapError(err, "failed to get OpenAI response")
	}

	if len(resp.Choices) == 0 {
		return nil, errors.NewAppError("no response from OpenAI")
	}

	// Get the translated content
	translatedContent := strings.TrimSpace(resp.Choices[0].Message.Content)

	// Parse the translated JSON back to map
	var translatedCurriculum map[string]interface{}
	if err := json.Unmarshal([]byte(translatedContent), &translatedCurriculum); err != nil {
		return nil, errors.WrapError(err, "failed to parse translated JSON")
	}

	return &dto.GenerateTranslationAIResponse{
		TranslatedCurriculum: translatedCurriculum,
	}, nil
}

// getLanguageName returns the full language name for the given code
func (uc *generateTranslationAIUseCase) getLanguageName(langCode string) string {
	switch langCode {
	case "pt":
		return "Portuguese"
	case "en":
		return "English"
	case "es":
		return "Spanish"
	default:
		return "Portuguese"
	}
}
