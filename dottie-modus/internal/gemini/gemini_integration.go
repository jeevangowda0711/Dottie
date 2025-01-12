package gemini

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type SymptomInput struct {
	Symptoms      []string `json:"symptoms"`
	CycleLength   int      `json:"cycle_length"`
	CycleDuration int      `json:"cycle_duration"`
	Age           int      `json:"age"`
}

type GeminiService struct{}

// AnalyzeSymptoms simulates advanced analysis via the Gemini service.
// Replace this stub with your actual logic.
func (g *GeminiService) AnalyzeSymptoms(input SymptomInput, contexts []string) (string, error) {
	// Implement the logic to analyze symptoms using Gemini
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey("Your API Key"))
	if err != nil {
		return "", errors.New("Issue in error")
	}
	model := client.GenerativeModel("gemini-1.5-flash-8b")
	model.SetTemperature(1)        // Controls randomness (0-1)
	model.SetTopK(40)              // Limits token sampling pool
	model.SetTopP(0.95)            // Nucleus sampling threshold
	model.SetMaxOutputTokens(8192) // Maximum response length
	listOfSymptoms := strings.Join(input.Symptoms, ",")

	prompt := "For symptoms like " + listOfSymptoms + "Generate 5 insights in terms of possible diagnosis,description of the diagnosis related to the symptom along with their probabilities ranging from low to high as a json response in a list of objects"
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatalf("Error generating content: %v", err)
	}
	fmt.Println(resp.Candidates[0].Content.Parts[0])
	var genaiText genai.Text
	var questions interface{}
	if len(resp.Candidates) > 0 {
		for _, part := range resp.Candidates[0].Content.Parts {
			switch p := part.(type) {
			case genai.Text:
				genaiText = p
				genaiText = genai.Text(strings.ReplaceAll(strings.ReplaceAll(string(genaiText), "```", ""), "json", ""))
				err := json.Unmarshal([]byte(genaiText), &questions)
				if err != nil {
					log.Fatalf("Error unmarshaling JSON: %v", err)
				}
			}
		}
	}
	val, _ := json.Marshal(questions)

	fmt.Println(string(val))
	return "Diagnosis from Gemini", nil
}

// IntegrateWithGemini is a convenience wrapper that instantiates GeminiService
// and delegates to AnalyzeSymptoms.
func IntegrateWithGemini(input SymptomInput, context []string) (string, error) {
	geminiService := &GeminiService{}
	return geminiService.AnalyzeSymptoms(input, context)
}
