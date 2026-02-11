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

// GenerateSkillAIUseCase defines the interface for AI filtering operations
type GenerateSkillAIUseCase interface {
	FilterContent(ctx context.Context, req *dto.GenerateSkillAIRequest) (*dto.GenerateSkillAIResponse, error)
}

// generateSkillAIUseCase implements GenerateSkillAIUseCase interface
type generateSkillAIUseCase struct {
	openaiClient *openai.Client
}

// NewGenerateSkillAIUseCase creates a new instance of GenerateSkillAIUseCase
func NewGenerateSkillAIUseCase(apiKey string) (GenerateSkillAIUseCase, error) {
	if apiKey == "" {
		return nil, errors.NewAppError("OPENAI_API_KEY environment variable is required")
	}

	client := openai.NewClient(option.WithAPIKey(apiKey))

	return &generateSkillAIUseCase{
		openaiClient: &client,
	}, nil
}

// FilterContent processes the content through OpenAI API to generate related skills
func (uc *generateSkillAIUseCase) FilterContent(ctx context.Context, req *dto.GenerateSkillAIRequest) (*dto.GenerateSkillAIResponse, error) {
	// Prepare the prompt for generating related skills
	prompt := fmt.Sprintf(`
Generate a list of up to 10 related skills for the following skill or area of expertise:

Skill/Area: %s

Provide only the list of related skills, without comments or additional analysis. The language must match the user's input language.
`, req.Content)

	// Create chat completion request
	chatReq := openai.ChatCompletionNewParams{
		Model: "gpt-4o-mini", // Using gpt-4o-mini
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(`You are a professional career advisor specialized in identifying and suggesting related skills that complement and enhance a person's professional profile.

Your task is to generate a list of up to 10 related skills based on a user-provided skill or area of expertise (e.g., Programming, Marketing, Design, Sales, Management, etc.).

BEFORE WRITING:
Analyze the input skill for clarity and coherence.

If the input is too vague, ask the user to provide a clearer skill or area of expertise.

Always interpret the input as the main skill or area, and derive a meaningful list of complementary and related skills that would be valuable in that field.

WHEN WRITING:
RANDOMLY choose one of the following tones (do not label or explain the tone):

Formal and concise

Dynamic and modern

Natural and conversational

Assertive and results-driven

Friendly and human (still professional)

Generate a unique list (max 10 items) of related and complementary skills that:

Are specifically related to the user's main skill or area of expertise

Include both technical and soft skills relevant to that field

Use professional, recruiter-friendly vocabulary

Vary in structure and tone (avoid rigid templates)

Reflect both hard skills and soft skills applicable to the field

Avoid buzzwords, overly technical jargon, and repeated structures

Use proper grammar, punctuation, and sentence flow

EXAMPLES OF WHAT TO INCLUDE:
- For Programming/Informatics: Problem Solving, Algorithm Design, Code Review, Version Control, Testing, Debugging, etc.
- For Marketing: Market Research, Content Creation, SEO, Social Media, Analytics, Brand Management, etc.
- For Design: User Experience, Color Theory, Typography, Prototyping, Adobe Creative Suite, etc.
- For Sales: Customer Relationship Management, Negotiation, Lead Generation, Presentation Skills, etc.

LANGUAGE RULES:
If the input is in Portuguese, return the list in Portuguese, using natural, correct, and professional vocabulary.

If the input is in English, return the list in English, following the same quality standards.

If the input is in Spanish, return the list in Spanish, using accurate and professional vocabulary.

Always match the language of the input.

ADDITIONAL INSTRUCTIONS:
Every new request must result in a new and varied list. Do not repeat formulas or templates.

If the input is unclear or insufficient, ask briefly for more detail (e.g., "Can you specify better the skill or area of expertise?").`),
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

	return &dto.GenerateSkillAIResponse{
		FilteredContent: filteredContent,
	}, nil
}
