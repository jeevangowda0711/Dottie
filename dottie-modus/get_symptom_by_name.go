package main

import (
	"fmt"
	"strconv"

	"github.com/hypermodeinc/modus/sdk/go/pkg/neo4j"
)

type Symptom struct {
	Name string `json:"name"`
}

type Condition struct {
	Name              string `json:"name"`
	Definition        string `json:"definition"`
	Severity          string `json:"severity"`
	RequiresAttention bool   `json:"requiresAttention"`
	Action            string `json:"action"`
}

type NormalRange struct {
	Name string `json:"name"`
	Min  int    `json:"min"`
	Max  int    `json:"max"`
	Unit string `json:"unit"`
}

// -----------------------------------------------------------------------------
// EXAMPLE QUERY #1: GetSymptomByName
// -----------------------------------------------------------------------------
func GetSymptomByName(name string) (Symptom, error) {
	query := "MATCH (s:Symptom {name: $name}) RETURN s.name AS name"
	params := map[string]interface{}{
		"name": name,
	}

	result, err := neo4j.ExecuteQuery("my-neo4j", query, params)
	if err != nil {
		return Symptom{}, fmt.Errorf("failed to execute query: %w", err)
	}

	if len(result.Records) == 0 {
		return Symptom{}, fmt.Errorf("no symptom found for name %q", name)
	}

	rawValue, ok := result.Records[0].Get("name")
	if !ok {
		return Symptom{}, fmt.Errorf("field 'name' not found in record")
	}

	// rawValue is a string already
	return Symptom{Name: rawValue}, nil
}

// -----------------------------------------------------------------------------
// INPUT/OUTPUT STRUCTS
// -----------------------------------------------------------------------------
type SymptomInput struct {
	CycleLength   int
	CycleDuration int
	Symptoms      []string
}

type AnalysisOutput struct {
	Diagnosis            string
	Recommendations      []string
	EducationalResources []string
}

// -----------------------------------------------------------------------------
// analyzeSymptoms
// -----------------------------------------------------------------------------
func analyzeSymptoms(input SymptomInput) (AnalysisOutput, error) {
	// 1) Check normal ranges
	normalRanges, err := queryNormalRanges(input.CycleLength, input.CycleDuration)
	if err != nil {
		return AnalysisOutput{}, err
	}
	if len(normalRanges) > 0 {
		// If we found normal ranges, mark as Normal
		return AnalysisOutput{
			Diagnosis:       "Normal",
			Recommendations: []string{"No action needed"},
		}, nil
	}

	// 2) Otherwise, diagnose abnormalities
	_, err = queryAbnormalities(input.CycleLength, input.CycleDuration)
	if err != nil {
		return AnalysisOutput{}, err
	}

	conditions, err := queryConditionsBySymptoms(input.Symptoms)
	if err != nil {
		return AnalysisOutput{}, err
	}

	_, err = queryCausesByConditions(conditions)
	if err != nil {
		return AnalysisOutput{}, err
	}

	// 3) Generate recommendations/resources
	recommendations := generateRecommendations(conditions)
	educationalResources := generateEducationalResources()

	return AnalysisOutput{
		Diagnosis:            "Abnormal",
		Recommendations:      recommendations,
		EducationalResources: educationalResources,
	}, nil
}

// -----------------------------------------------------------------------------
// QUERY #1: Normal Ranges
// -----------------------------------------------------------------------------
func queryNormalRanges(cycleLength, cycleDuration int) ([]NormalRange, error) {
	query := `
        MATCH (n:NormalRange)
        RETURN n.name AS name, n.min AS min, n.max AS max, n.unit AS unit
    `
	params := map[string]interface{}{
		"cycleLength":   cycleLength,
		"cycleDuration": cycleDuration,
	}

	result, err := neo4j.ExecuteQuery("my-neo4j", query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	var normalRanges []NormalRange
	for _, record := range result.Records {
		// 1) name
		nameVal, ok := record.Get("name")
		if !ok {
			return nil, fmt.Errorf("field 'name' not found")
		}
		// nameVal is a string
		nameStr := nameVal

		// 2) min
		minVal, ok := record.Get("min")
		if !ok {
			return nil, fmt.Errorf("field 'min' not found")
		}
		// parse the string to float
		floatMin, err := strconv.ParseFloat(minVal, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse 'min'=%q as float64: %v", minVal, err)
		}

		// 3) max
		maxVal, ok := record.Get("max")
		if !ok {
			return nil, fmt.Errorf("field 'max' not found")
		}
		floatMax, err := strconv.ParseFloat(maxVal, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse 'max'=%q as float64: %v", maxVal, err)
		}

		// 4) unit
		unitVal, ok := record.Get("unit")
		if !ok {
			return nil, fmt.Errorf("field 'unit' not found")
		}
		unitStr := unitVal

		normalRanges = append(normalRanges, NormalRange{
			Name: nameStr,
			Min:  int(floatMin),
			Max:  int(floatMax),
			Unit: unitStr,
		})
	}

	return normalRanges, nil
}

// -----------------------------------------------------------------------------
// QUERY #2: Abnormalities
// -----------------------------------------------------------------------------
func queryAbnormalities(cycleLength, cycleDuration int) ([]string, error) {
	query := `
        MATCH (a:Abnormality)
        RETURN a.description AS description
    `
	params := map[string]interface{}{
		"cycleLength":   cycleLength,
		"cycleDuration": cycleDuration,
	}

	result, err := neo4j.ExecuteQuery("my-neo4j", query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	var abnormalities []string
	for _, record := range result.Records {
		descVal, ok := record.Get("description")
		if !ok {
			return nil, fmt.Errorf("field 'description' not found")
		}
		// descVal is a string
		abnormalities = append(abnormalities, descVal)
	}

	return abnormalities, nil
}

// -----------------------------------------------------------------------------
// QUERY #3: Conditions by Symptoms
// -----------------------------------------------------------------------------
func queryConditionsBySymptoms(symptoms []string) ([]Condition, error) {
	query := `
        MATCH (c:Condition)-[:CAUSES]->(s:Symptom)
        WHERE s.name IN $symptoms
        RETURN c.name AS name, c.definition AS definition, c.severity AS severity,
               c.requiresAttention AS requiresAttention, c.action AS action
    `
	params := map[string]interface{}{
		"symptoms": symptoms,
	}

	result, err := neo4j.ExecuteQuery("my-neo4j", query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	var conditions []Condition
	for _, record := range result.Records {
		// c.name
		nameVal, ok := record.Get("name")
		if !ok {
			return nil, fmt.Errorf("field 'name' not found in Condition record")
		}

		// c.definition
		defVal, ok := record.Get("definition")
		if !ok {
			return nil, fmt.Errorf("field 'definition' not found")
		}

		// c.severity
		sevVal, ok := record.Get("severity")
		if !ok {
			return nil, fmt.Errorf("field 'severity' not found")
		}

		// c.requiresAttention (string "true"/"false"?)
		attVal, ok := record.Get("requiresAttention")
		if !ok {
			return nil, fmt.Errorf("field 'requiresAttention' not found")
		}
		boolVal, err := strconv.ParseBool(attVal)
		if err != nil {
			return nil, fmt.Errorf("cannot parse requiresAttention=%q as bool: %v", attVal, err)
		}

		// c.action
		actionVal, ok := record.Get("action")
		if !ok {
			return nil, fmt.Errorf("field 'action' not found")
		}

		conditions = append(conditions, Condition{
			Name:              nameVal,
			Definition:        defVal,
			Severity:          sevVal,
			RequiresAttention: boolVal,
			Action:            actionVal,
		})
	}

	return conditions, nil
}

// -----------------------------------------------------------------------------
// QUERY #4: Causes by Conditions
// -----------------------------------------------------------------------------
func queryCausesByConditions(conditions []Condition) ([]string, error) {
	query := `
        MATCH (c:Condition)-[:CAUSES]->(cause:Cause)
        WHERE c.name IN $conditions
        RETURN cause.name AS name
    `
	// Convert Condition structs -> list of names
	var conditionNames []string
	for _, cond := range conditions {
		conditionNames = append(conditionNames, cond.Name)
	}

	params := map[string]interface{}{
		"conditions": conditionNames,
	}

	result, err := neo4j.ExecuteQuery("my-neo4j", query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	var causes []string
	for _, record := range result.Records {
		rawName, ok := record.Get("name")
		if !ok {
			return nil, fmt.Errorf("field 'name' not found in Cause")
		}
		// rawName is a string
		causes = append(causes, rawName)
	}

	return causes, nil
}

// -----------------------------------------------------------------------------
// GENERATION FUNCTIONS
// -----------------------------------------------------------------------------
func generateRecommendations(conditions []Condition) []string {
	recommendations := []string{}
	for _, condition := range conditions {
		if condition.Severity == "high" {
			recommendations = append(recommendations, "Seek medical attention immediately.")
		} else {
			recommendations = append(recommendations, "Monitor and consult a doctor if it persists.")
		}
	}
	return recommendations
}

func generateEducationalResources() []string {
	// Stubbed-out logic
	return []string{}
}

// -----------------------------------------------------------------------------
// MAIN (EXAMPLE USAGE)
// -----------------------------------------------------------------------------
func main() {
	symptom, err := GetSymptomByName("Dysmenorrhea")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Symptom:", symptom.Name)

	input := SymptomInput{
		CycleLength:   28,
		CycleDuration: 5,
		Symptoms:      []string{"Dysmenorrhea"},
	}
	analysis, err := analyzeSymptoms(input)
	if err != nil {
		fmt.Println("Error in analysis:", err)
		return
	}
	fmt.Println("Diagnosis:", analysis.Diagnosis)
	fmt.Println("Recommendations:", analysis.Recommendations)
	fmt.Println("Educational Resources:", analysis.EducationalResources)
}
