package main

import (
	"fmt"

	"github.com/hypermodeinc/modus/sdk/go/pkg/neo4j"
)

type Symptom struct {
	Name string `json:"name"`
}

func GetSymptomByName(name string) (Symptom, error) {
	query := "MATCH (s:Symptom {name: $name}) RETURN s.name AS name"
	params := map[string]interface{}{
		"name": name,
	}

	// Use Modus SDK to execute Neo4j query
	result, err := neo4j.ExecuteQuery("my-neo4j", query, params)
	if err != nil {
		return Symptom{}, fmt.Errorf("failed to execute query: %w", err)
	}

	if len(result.Records) == 0 {
		return Symptom{}, fmt.Errorf("no records found")
	}

	// Correct handling of result.Records[0].Get
	value, exists := result.Records[0].Get("name")
	if !exists {
		return Symptom{}, fmt.Errorf("failed to get symptom name: field 'name' not found")
	}

	return Symptom{Name: value}, nil
}

func main() {
	// Example usage
	symptom, err := GetSymptomByName("Cough")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Symptom:", symptom.Name)
	}
}
