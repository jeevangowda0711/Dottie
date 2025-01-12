// data_population.go
// This file populates the Neo4j database using the extracted ACOG data from a JSON file.

package db

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hypermodeinc/modus/sdk/go/pkg/neo4j"
)

// LoadACOGData reads the JSON file, parses the data, and populates the Neo4j database
func LoadACOGData(filePath string) error {
	fmt.Printf("Loading data from: %s\n", filePath)

	// Open the JSON file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Decode the JSON file
	var data ACOGData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	fmt.Printf("Found %d normal ranges, %d conditions, %d symptoms, and %d educational content\n",
		len(data.NormalRanges), len(data.Conditions), len(data.Symptoms), len(data.EducationalContent))

	// Create nodes in Neo4j
	for _, nr := range data.NormalRanges {
		fmt.Printf("Creating normal range node: %s\n", nr.Name)
		if err := createNormalRangeNode(nr); err != nil {
			return fmt.Errorf("failed to create normal range node %s: %w", nr.Name, err)
		}
	}

	for _, condition := range data.Conditions {
		fmt.Printf("Creating condition node: %s\n", condition.Name)
		if err := createConditionNode(condition); err != nil {
			return fmt.Errorf("failed to create condition node %s: %w", condition.Name, err)
		}
	}

	for _, symptom := range data.Symptoms {
		fmt.Printf("Creating symptom node: %s\n", symptom.Name)
		if err := createSymptomNode(symptom); err != nil {
			return fmt.Errorf("failed to create symptom node %s: %w", symptom.Name, err)
		}
	}

	for _, content := range data.EducationalContent {
		fmt.Printf("Creating educational content node: %s\n", content.Title)
		if err := createEducationalContentNode(content); err != nil {
			return fmt.Errorf("failed to create educational content node %s: %w", content.Title, err)
		}
	}

	fmt.Println("Database populated successfully!")
	return nil
}

func createNormalRangeNode(nr NormalRange) error {
	query := `
	CREATE (n:NormalRange {
		name: $name,
		min: $min,
		max: $max,
		unit: $unit
	}) RETURN n`

	params := map[string]interface{}{
		"name": nr.Name,
		"min":  nr.Min,
		"max":  nr.Max,
		"unit": nr.Unit,
	}

	result, err := neo4j.ExecuteQuery("my-neo4j", query, params)
	if err != nil {
		return fmt.Errorf("error creating normal range node: %w", err)
	}

	fmt.Printf("Created normal range node: %s, result: %+v\n", nr.Name, result)
	return nil
}

func createConditionNode(condition Condition) error {
	query := `
	CREATE (c:Condition {
		name: $name,
		definition: $definition,
		severity: $severity,
		requiresAttention: $requiresAttention
	}) RETURN c`

	params := map[string]interface{}{
		"name":              condition.Name,
		"definition":        condition.Definition,
		"severity":          condition.Severity,
		"requiresAttention": condition.RequiresAttention,
	}

	result, err := neo4j.ExecuteQuery("my-neo4j", query, params)
	if err != nil {
		return fmt.Errorf("error creating condition node: %w", err)
	}

	fmt.Printf("Created condition node: %s, result: %+v\n", condition.Name, result)
	return nil
}

func createSymptomNode(symptom Symptom) error {
	query := `
	CREATE (s:Symptom {
		name: $name,
		description: $description,
		severity: $severity
	}) RETURN s`

	params := map[string]interface{}{
		"name":        symptom.Name,
		"description": symptom.Description,
		"severity":    symptom.Severity,
	}

	result, err := neo4j.ExecuteQuery("my-neo4j", query, params)
	if err != nil {
		return fmt.Errorf("error creating symptom node: %w", err)
	}

	fmt.Printf("Created symptom node: %s, result: %+v\n", symptom.Name, result)
	return nil
}

func createEducationalContentNode(content EducationalContent) error {
	query := `
	CREATE (e:EducationalContent {
		type: $type,
		url: $url,
		title: $title,
		source: $source
	}) RETURN e`

	params := map[string]interface{}{
		"type":   content.Type,
		"url":    content.Url,
		"title":  content.Title,
		"source": content.Source,
	}

	result, err := neo4j.ExecuteQuery("my-neo4j", query, params)
	if err != nil {
		return fmt.Errorf("error creating educational content node: %w", err)
	}

	fmt.Printf("Created educational content node: %s, result: %+v\n", content.Title, result)
	return nil
}
