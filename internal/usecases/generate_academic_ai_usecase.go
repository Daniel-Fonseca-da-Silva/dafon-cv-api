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

// GenerateAcademicAIUseCase defines the interface for AI filtering operations
type GenerateAcademicAIUseCase interface {
	FilterContent(ctx context.Context, req *dto.GenerateAcademicAIRequest) (*dto.GenerateAcademicAIResponse, error)
}

// generateAcademicAIUseCase implements GenerateAcademicAIUseCase interface
type generateAcademicAIUseCase struct {
	openaiClient *openai.Client
}

// NewGenerateAcademicAIUseCase creates a new instance of GenerateAcademicAIUseCase
func NewGenerateAcademicAIUseCase(apiKey string) (GenerateAcademicAIUseCase, error) {
	if apiKey == "" {
		return nil, errors.NewAppError("OPENAI_API_KEY environment variable is required")
	}

	client := openai.NewClient(option.WithAPIKey(apiKey))

	return &generateAcademicAIUseCase{
		openaiClient: &client,
	}, nil
}

// FilterContent processes the content through OpenAI API to filter and improve it
func (uc *generateAcademicAIUseCase) FilterContent(ctx context.Context, req *dto.GenerateAcademicAIRequest) (*dto.GenerateAcademicAIResponse, error) {
	// Prepare the prompt for filtering
	prompt := fmt.Sprintf(`
Generate a list of relevant academic activities, subjects, and experiences for the following university degree or field of study:

Degree/Field: %s

Provide only the list of academic activities, without comments or additional analysis. The language must match the user's input language.
`, req.Content)

	// Create chat completion request
	chatReq := openai.ChatCompletionNewParams{
		Model: "gpt-4o-mini", // Using gpt-4o-mini
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(`You are a professional resume writer specialized in crafting impactful and recruiter-friendly academic activity lists that enhance resumes.

Your task is to generate a list of relevant academic activities, subjects, and experiences based on a user-provided university degree or field of study (e.g., Computer Science, Medicine, Engineering, Business Administration, etc.).

BEFORE WRITING:
Analyze the input for clarity and coherence.

If the input is too vague, ask the user to provide a clearer degree or field of study.

Always interpret the input as the main academic field, and derive a meaningful list of relevant subjects, activities, projects, or experiences typical for that degree.

WHEN WRITING:
RANDOMLY choose one of the following tones (do not label or explain the tone):

Formal and concise

Dynamic and modern

Natural and conversational

Assertive and results-driven

Friendly and human (still professional)

Generate a unique list (min 10 items and max 20 items) of relevant and realistic academic activities that:

Are specifically tailored to the user's degree/field of study

Include core subjects, practical activities, projects, and experiences typical for that academic area

Use professional, recruiter-friendly vocabulary

Vary in structure and tone (avoid rigid templates)

Reflect both theoretical knowledge and practical skills applicable to the field

Avoid buzzwords, overly technical jargon, and repeated structures

Use proper grammar, punctuation, and sentence flow

EXAMPLES OF WHAT TO INCLUDE:
- For Computer Science: Data Structures, Programming Logic, Software Engineering, Database Systems, etc.
- For Medicine: Human Anatomy, Physiology, Pathology, Clinical Practice, Medical Ethics, etc.
- For Engineering: Mathematics, Physics, Technical Drawing, Project Management, etc.
- For Business: Economics, Marketing, Management, Finance, Strategic Planning, etc.

LANGUAGE RULES:
If the input is in Portuguese, return the list in Portuguese, using natural, correct, and professional vocabulary.

If the input is in English, return the list in English, following the same quality standards.

If the input is in Spanish, return the list in Spanish, using accurate and professional vocabulary.

Always match the language of the input.

ADDITIONAL INSTRUCTIONS:
Every new request must result in a new and varied list. Do not repeat formulas or templates.

If the input is unclear or insufficient, ask briefly for more detail (e.g., "Can you specify better the degree or field of study?").`),
			openai.UserMessage(prompt),
		},
		MaxTokens:   openai.Int(int64(config.ParseIntEnv("OPENAI_MAX_TOKENS", 1000))),
		Temperature: openai.Float(config.ParseFloatEnv("OPENAI_TEMPERATURE", 0.7)),
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

	// Get the filtered content
	filteredContent := resp.Choices[0].Message.Content

	return &dto.GenerateAcademicAIResponse{
		FilteredContent: filteredContent,
	}, nil
}
