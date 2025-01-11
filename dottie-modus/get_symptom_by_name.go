//go:build !js
// +build !js

package main

import (
	"context"
	"fmt"

	"github.com/hypermodeinc/modus/sdk/go/pkg/neo4j"
)

type Symptom struct {
	Name string `json:"name"`
}

func GetSymptomByName(ctx context.Context, name string) (Symptom, error) {
	query := "MATCH (s:Symptom {name: $name}) RETURN s.name as name"
	params := map[string]any{
		"name": name,
	}

	// Use Modus SDK to execute Neo4j query
	result, err := neo4j.ExecuteQuery("my-neo4j", query, params)
	if err != nil {
		return Symptom{}, fmt.Errorf("failed to execute query: %w", err)
	}

	if len(result.Records) == 0 {
		return Symptom{}, fmt.Errorf("symptom not found")
	}

	symptomName, err := result.Records[0].Get("name")
	if err != nil {
		return Symptom{}, fmt.Errorf("failed to get symptom name: %w", err)
	}

	return Symptom{Name: symptomName}, nil
}

func main() {
	// Example usage
	ctx := context.Background()
	symptom, err := GetSymptomByName(ctx, "Cough")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Symptom:", symptom.Name)
	}
}
