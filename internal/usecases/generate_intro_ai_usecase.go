package usecases

import (
	"context"
	"fmt"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/config"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/errors"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// GenerateIntroAIUseCase defines the interface for AI filtering operations
type GenerateIntroAIUseCase interface {
	FilterContent(ctx context.Context, req *dto.GenerateIntroAIRequest) (*dto.GenerateIntroAIResponse, error)
}

// generateIntroAIUseCase implements GenerateIntroAIUseCase interface
type generateIntroAIUseCase struct {
	openaiClient *openai.Client
}

// NewGenerateIntroAIUseCase creates a new instance of GenerateIntroAIUseCase
func NewGenerateIntroAIUseCase(apiKey string) (GenerateIntroAIUseCase, error) {
	if apiKey == "" {
		return nil, errors.NewAppError("OPENAI_API_KEY environment variable is required")
	}

	client := openai.NewClient(option.WithAPIKey(apiKey))

	return &generateIntroAIUseCase{
		openaiClient: &client,
	}, nil
}

// FilterContent processes the content through OpenAI API to filter and improve it
func (uc *generateIntroAIUseCase) FilterContent(ctx context.Context, req *dto.GenerateIntroAIRequest) (*dto.GenerateIntroAIResponse, error) {
	// Prepare the prompt for filtering
	prompt := fmt.Sprintf(`
Improve and filter the following content, making it more professional and appropriate for a resume:

Content: %s

Provide only the improved description, without comments or additional analysis, the language must be following the user's language.
`, req.Content)

	// Create chat completion request
	chatReq := openai.ChatCompletionNewParams{
		Model: "gpt-4o-mini", // Using gpt-4o-mini
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(`You are a professional resume writer specialized in crafting impactful and concise professional summaries that enhance resumes.

Your task is to generate a resume-ready description (max. 320 characters) based on user input. This description should be unique, varied in tone, grammatically correct, and recruiter-friendly.

BEFORE WRITING:
- Analyze the input for clarity, coherence, and completeness.
- Validate if the information is sufficient to create a meaningful description. If it's too vague or lacks accomplishments, highlight what's missing (but don't generate the description).
- If multiple roles/skills are listed, prioritize the most impactful or recent ones.

WHEN WRITING:
1. RANDOMLY choose one of the following tones:
   - Formal and concise
   - Dynamic and modern
   - Natural and conversational
   - Assertive and results-driven
   - Friendly and human (still professional)

2. Generate a 1–2 sentence description (max 320 characters) that:
   - Highlights achievements, results, or measurable impact (when possible)
   - Uses past-tense action verbs
   - Avoids buzzwords, excessive adjectives, or technical jargon
   - Is clear, objective, and suitable for resume use
   - Has a natural sentence structure, not robotic or repetitive

3. DO NOT label or explain the tone used. Only return the final description.

ADDITIONAL INSTRUCTIONS:
- Each new description request must generate a completely new sentence structure and style.
- Avoid repeating templates or formulas from previous outputs.
- Always tailor the writing to make the person stand out positively — even in short form.

IF INPUT IS IN PORTUGUESE:
- Return the description in Portuguese, following the same structure and tone logic.
- Use natural, correct Portuguese with professional vocabulary.

IF INPUT IS IN ENGLISH:
- Return the description in English, following the same structure and tone logic.
- Use natural, correct English with professional vocabulary.

IF INPUT IS IN SPANISH:
- Return the description in Spanish, following the same structure and tone logic.
- Use natural, correct Spanish with professional vocabulary.

If the information is insufficient or unclear:
- Reply with a short message asking for more detail (e.g., specific achievements, role type, or impact).`),
			openai.UserMessage(prompt),
		},
		MaxTokens:   openai.Int(int64(config.ParseIntEnv("OPENAI_MAX_TOKENS", 1000))),
		Temperature: openai.Float(config.ParseFloatEnv("OPENAI_TEMPERATURE", 0.7)),
		TopP:        openai.Float(config.ParseFloatEnv("OPENAI_TOP_P", 1.0)),
	}

	// Call OpenAI API
	resp, err := uc.openaiClient.Chat.Completions.New(ctx, chatReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get OpenAI response for intro generation: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI for intro generation")
	}

	// Get the filtered content
	filteredContent := resp.Choices[0].Message.Content

	return &dto.GenerateIntroAIResponse{
		FilteredContent: filteredContent,
	}, nil
}
