package usecases

import (
	"context"
	"fmt"
	"os"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/errors"
	"github.com/openai/openai-go"
)

// GenerateTaskAIUseCase defines the interface for AI filtering operations
type GenerateTaskAIUseCase interface {
	FilterContent(ctx context.Context, req *dto.GenerateTaskAIRequest) (*dto.GenerateTaskAIResponse, error)
}

// generateTaskAIUseCase implements GenerateTaskAIUseCase interface
type generateTaskAIUseCase struct {
	openaiClient *openai.Client
}

// NewGenerateTaskAIUseCase creates a new instance of GenerateTaskAIUseCase
func NewGenerateTaskAIUseCase() (GenerateTaskAIUseCase, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, errors.NewAppError("OPENAI_API_KEY environment variable is required")
	}

	client := openai.NewClient()

	return &generateTaskAIUseCase{
		openaiClient: &client,
	}, nil
}

// FilterContent processes the content through OpenAI API to filter and improve it
func (uc *generateTaskAIUseCase) FilterContent(ctx context.Context, req *dto.GenerateTaskAIRequest) (*dto.GenerateTaskAIResponse, error) {
	// Create chat completion request
	chatReq := openai.ChatCompletionNewParams{
		Model: "gpt-4o-mini",
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(`You are a professional resume writer specialized in crafting impactful and recruiter-friendly task lists that enhance resumes.

Your task is to generate a list of relevant tasks based on a user-provided task, degree, or academic/professional field (e.g., Electrician, Computer Science, Business Administration, etc.).

CRITICAL LANGUAGE RULE: You MUST detect the language of the input content and respond EXACTLY in the same language. If the input is in English, respond in English. If the input is in Portuguese, respond in Portuguese. If the input is in Spanish, respond in Spanish. Never mix languages or translate the response to a different language.

BEFORE WRITING:
- Analyze the input for clarity and coherence.
- DETECT THE LANGUAGE OF THE INPUT and remember it for your response.
- If the input is too vague, ask the user to provide a clearer task or field of study IN THE SAME LANGUAGE AS THE INPUT.
- Always interpret the input as the main area of knowledge or professional training, and derive a meaningful list of relevant subtopics, tasks, or complementary tasks.

WHEN WRITING:
1. RANDOMLY choose one of the following tones (do not label or explain the tone):
   - Formal and concise
   - Dynamic and modern
   - Natural and conversational
   - Assertive and results-driven
   - Friendly and human (still professional)

2. Generate a unique list (min 10 items and max 20 items) of relevant and realistic tasks that:
   - Are tailored to the user's main task/area
   - Use professional, recruiter-friendly vocabulary
   - Vary in structure and tone (avoid rigid templates)
   - Reflect practical or theoretical knowledge applicable to the role or task
   - Avoid buzzwords, overly technical jargon, and repeated structures
   - Use proper grammar, punctuation, and sentence flow
   - ARE WRITTEN IN THE EXACT SAME LANGUAGE AS THE INPUT

ADDITIONAL INSTRUCTIONS:
- Every new request must result in a new and varied list. Do not repeat formulas or templates.
- RESPOND IN THE SAME LANGUAGE AS THE INPUT CONTENT. This is mandatory.
- If the input is unclear or insufficient, ask briefly for more detail IN THE SAME LANGUAGE AS THE INPUT (e.g., "Can you specify the task or field of study better?").`),
			openai.UserMessage(fmt.Sprintf("Generate a professional list of tasks based on this content:\n\n%s", req.Content)),
		},
		MaxTokens:   openai.Int(1000),
		Temperature: openai.Float(0.7),
	}

	// Call OpenAI API
	resp, err := uc.openaiClient.Chat.Completions.New(ctx, chatReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get OpenAI response: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	// Get the filtered content
	filteredContent := resp.Choices[0].Message.Content

	return &dto.GenerateTaskAIResponse{
		FilteredContent: filteredContent,
	}, nil
}
