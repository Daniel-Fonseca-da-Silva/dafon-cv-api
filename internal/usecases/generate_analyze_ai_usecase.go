package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/dto"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/errors"
	"github.com/google/uuid"
	"github.com/openai/openai-go"
)

// GenerateAnalyzeAIUseCase defines the interface for AI filtering operations
type GenerateAnalyzeAIUseCase interface {
	FilterContent(ctx context.Context, curriculumID uuid.UUID, language string) (*dto.GenerateAnalyzeAIResponse, error)
}

// generateAnalyzeAIUseCase implements GenerateAnalyzeAIUseCase interface
type generateAnalyzeAIUseCase struct {
	openaiClient      *openai.Client
	curriculumUseCase CurriculumUseCase
}

// NewGenerateAnalyzeAIUseCase creates a new instance of GenerateAnalyzeAIUseCase
func NewGenerateAnalyzeAIUseCase(curriculumUseCase CurriculumUseCase) (GenerateAnalyzeAIUseCase, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, errors.NewAppError("OPENAI_API_KEY environment variable is required")
	}

	client := openai.NewClient()

	return &generateAnalyzeAIUseCase{
		openaiClient:      &client,
		curriculumUseCase: curriculumUseCase,
	}, nil
}

func (uc *generateAnalyzeAIUseCase) FilterContent(ctx context.Context, curriculumID uuid.UUID, language string) (*dto.GenerateAnalyzeAIResponse, error) {
	// Get curriculum body using the existing method from CurriculumUseCase
	curriculumBody, err := uc.curriculumUseCase.GetCurriculumBody(ctx, curriculumID)
	if err != nil {
		return nil, errors.WrapError(err, "failed to get curriculum body")
	}

	// Map language code to full language name for the prompt
	languageMap := map[string]string{
		"pt": "português",
		"en": "english",
		"es": "español",
	}
	languageName := languageMap[language]
	if languageName == "" {
		languageName = "english" // default fallback
	}

	// Prepare the prompt for curriculum analysis using the fetched body
	prompt := fmt.Sprintf(`
Analyze the following curriculum text and provide a comprehensive professional analysis. The curriculum contains personal information, experience, skills, and academic background.

Curriculum Text: %s

IMPORTANT: You MUST respond in %s language. All fields in the JSON response must be in %s.

Return ONLY a strict JSON object in %s language, with this exact structure and keys:
{
  "score": number, // 0-100 numeric score (can be decimal like 75.5)
  "description": string, // brief professional summary of the overall assessment
  "improvement_points": [string],
  "best_practices": [string],
  "ats_compatibility": {
    "assessment": string, // short assessment of ATS readiness
    "chance": string, // qualitative chance, e.g., "low", "medium", "high"
    "recommendations": [string]
  },
  "professional_alignment": [string],
  "strengths": [string],
  "recommendations": [string]
}

STRICT RULES:
- Output must be valid JSON only (no markdown, no backticks, no extra text)
- Keep responses concise and professional
- ALL content must be in %s language
`, curriculumBody.Body, languageName, languageName, languageName, languageName)

	// Create chat completion request
	chatReq := openai.ChatCompletionNewParams{
		Model: "gpt-4o-mini", // Using gpt-4o-mini
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(fmt.Sprintf(`You are a professional career consultant and resume expert with extensive experience in HR, recruitment, and career development. You specialize in analyzing curricula and providing comprehensive feedback to help professionals improve their job application success.

Your expertise includes:
- Modern resume best practices and industry standards
- ATS (Applicant Tracking System) optimization
- Professional sector alignment and career path analysis
- International recruitment standards and expectations
- Career development and professional growth strategies

ANALYSIS FRAMEWORK:
When analyzing a curriculum, you should:

1. **Overall Assessment**: Provide a numerical score (0-100) based on:
   - Content completeness and organization
   - Professional presentation
   - Relevance and impact of information
   - Grammar and language quality
   - Industry standards compliance

2. **Improvement Points**: Identify specific weaknesses such as:
   - Missing information or sections
   - Poor formatting or structure
   - Weak action verbs or descriptions
   - Lack of quantifiable achievements
   - Inappropriate content or tone

3. **Best Practices**: Recommend current industry standards:
   - Modern formatting and layout
   - Effective use of keywords
   - Professional language and tone
   - Appropriate length and structure
   - Industry-specific requirements

4. **ATS Compatibility**: Assess:
   - Keyword optimization
   - Format compatibility
   - Section organization
   - Technical requirements
   - Chances of passing automated screening

5. **Professional Alignment**: Evaluate:
   - Suitability for specific sectors (medicine, technology, manufacturing, logistics, etc.)
   - Career path coherence
   - Industry-specific requirements
   - Professional development trajectory

6. **Strengths**: Highlight:
   - Strong achievements and experiences
   - Relevant skills and qualifications
   - Professional growth indicators
   - Unique value propositions

7. **Actionable Recommendations**: Provide specific, implementable steps for improvement.

LANGUAGE AND FORMATTING REQUIREMENTS:
- You MUST respond in %s language
- Use professional, clear, and constructive language
- Maintain a supportive and encouraging tone
- Provide specific, actionable advice
- Be honest but diplomatic in your assessment
- Respond with strict JSON only, matching the requested schema
- Do not include markdown, code fences, or additional text outside JSON
- ALL fields, descriptions, and recommendations must be in %s language

Your goal is to help the person improve their curriculum and increase their chances of success in the job market.`, languageName, languageName)),
			openai.UserMessage(prompt),
		},
		MaxTokens:   openai.Int(2000),
		Temperature: openai.Float(0.7),
	}

	// Call OpenAI API
	resp, err := uc.openaiClient.Chat.Completions.New(ctx, chatReq)
	if err != nil {
		return nil, errors.WrapError(err, "failed to get OpenAI response")
	}

	if len(resp.Choices) == 0 {
		return nil, errors.NewAppError("no response from OpenAI")
	}

	// Get the analyzed content
	analyzedContent := resp.Choices[0].Message.Content

	// Attempt to unmarshal strict JSON into response DTO
	var structured dto.GenerateAnalyzeAIResponse
	if err := json.Unmarshal([]byte(analyzedContent), &structured); err != nil {
		// If parsing fails, return error
		return nil, errors.WrapError(err, "failed to parse AI response as JSON")
	}

	return &structured, nil
}