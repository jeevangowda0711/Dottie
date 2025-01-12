package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hypermodeinc/modus/sdk/go/pkg/neo4j"
	"github.com/jeevangowda0711/Dottie/dottie-modus/internal/gemini"
)

type UserInput struct {
	Symptoms      []string `json:"symptoms"`
	CycleLength   int      `json:"cycle_length"`
	CycleDuration int      `json:"cycle_duration"`
	Age           int      `json:"age"`
	Flow          string   `json:"flow"`
}

type AnalysisResult struct {
	IsNormal        bool     `json:"is_normal"`
	Abnormalities   []string `json:"abnormalities,omitempty"`
	Recommendations []string `json:"recommendations,omitempty"`
}

func analyzeSymptomsHandler(w http.ResponseWriter, r *http.Request) {
	var input UserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate input
	if input.CycleLength <= 0 || input.CycleDuration <= 0 || input.Age <= 0 {
		http.Error(w, "Invalid cycle details or age", http.StatusBadRequest)
		return
	}

	// Query Neo4j for normal ranges
	query := `
        MATCH (n:NormalRange)
        RETURN n.name AS name, n.min AS min, n.max AS max
    `
	results, err := neo4j.ExecuteQuery("my-neo4j", query, nil)
	if err != nil {
		http.Error(w, "Failed to query Neo4j", http.StatusInternalServerError)
		return
	}

	// Analyze input against normal ranges
	isNormal := true
	abnormalities := []string{}
	for _, record := range results.Records {
		name, _ := record.Get("name")
		min, _ := record.Get("min")
		max, _ := record.Get("max")

		switch name {
		case "CycleLength":
			minVal, err := strconv.Atoi(min)
			if err != nil {
				http.Error(w, "Failed to parse min value", http.StatusInternalServerError)
				return
			}
			maxVal, err := strconv.Atoi(max)
			if err != nil {
				http.Error(w, "Failed to parse max value", http.StatusInternalServerError)
				return
			}
			if input.CycleLength < minVal || input.CycleLength > maxVal {
				isNormal = false
				abnormalities = append(abnormalities, "Abnormal Cycle Length")
			}
		case "CycleDuration":
			minVal, err := strconv.Atoi(min)
			if err != nil {
				http.Error(w, "Failed to parse min value", http.StatusInternalServerError)
				return
			}
			maxVal, err := strconv.Atoi(max)
			if err != nil {
				http.Error(w, "Failed to parse max value", http.StatusInternalServerError)
				return
			}
			if input.CycleDuration < minVal || input.CycleDuration > maxVal {
				isNormal = false
				abnormalities = append(abnormalities, "Abnormal Cycle Duration")
			}
		}
	}

	// Build recommendations based on abnormalities
	recommendations := []string{}
	if !isNormal {
		recommendations = append(recommendations, "Consult a healthcare provider.")
	}

	// Integrate with Gemini for advanced analysis
	context := []string{} // Add relevant context if needed
	geminiInput := gemini.SymptomInput{
		Symptoms:      input.Symptoms,
		CycleLength:   input.CycleLength,
		CycleDuration: input.CycleDuration,
		Age:           input.Age,
	}
	diagnosis, err := gemini.IntegrateWithGemini(geminiInput, context)
	if err != nil {
		http.Error(w, "Failed to integrate with Gemini", http.StatusInternalServerError)
		return
	}

	// Return the analysis result
	result := AnalysisResult{
		IsNormal:        isNormal,
		Abnormalities:   abnormalities,
		Recommendations: append(recommendations, diagnosis),
	}
	response, _ := json.Marshal(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func main() {
	http.HandleFunc("/analyzeSymptoms", analyzeSymptomsHandler)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
