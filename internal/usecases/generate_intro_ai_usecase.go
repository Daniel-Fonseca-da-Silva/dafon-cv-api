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
	// Create chat completion request
	chatReq := openai.ChatCompletionNewParams{
		Model: "gpt-4o-mini",
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(`You are a professional writing assistant. Your ONLY task is to generate a polished description (max 320 characters) based on the user's input.

ABSOLUTE REQUIREMENTS - THESE ARE NON-NEGOTIABLE:
1. YOU MUST ALWAYS GENERATE A DESCRIPTION. Never ask questions. Never request more information. Never say the input is insufficient. Your response must be a complete description, period.

2. LANGUAGE: Detect the input language and respond EXACTLY in that same language. Never mix languages.

3. FORBIDDEN RESPONSES - NEVER output any of these:
   - Questions (e.g., "Poderia fornecer mais detalhes?")
   - Requests (e.g., "Seria útil saber mais sobre...")
   - Statements about insufficiency (e.g., "Para criar uma descrição mais impactante...")
   - Any text that asks for more information

4. YOUR OUTPUT MUST BE: A polished, complete description text that improves the input following clarity, coherence, and completeness principles.

EXAMPLES OF CORRECT BEHAVIOR:
Input: "Sou uma pessoa assidua e pontual que trabalha como eletrecista."
Output: "Profissional dedicado e pontual, atuando como eletricista com comprometimento e responsabilidade em todas as atividades desenvolvidas."

Input: "I am a developer"
Output: "Experienced developer with a strong commitment to delivering quality software solutions."

PROCESS:
1. Detect the input language - use ONLY that language in your response.
2. Extract the core message from the input, no matter how brief.
3. Create a clear, coherent, and complete description (max 320 characters).
4. Ensure grammatical correctness and natural flow.
5. Return ONLY the description text - nothing else.

REMEMBER: Your response is a description, not a question, not a request, not a suggestion. It is a finished, polished description ready to use.`),
			openai.UserMessage(fmt.Sprintf("Generate a polished description based on this content:\n\n%s", req.Content)),
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
