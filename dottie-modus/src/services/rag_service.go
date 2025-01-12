package services

import (
	"context"
	"fmt"

	"github.com/hypermodeinc/modus/sdk/go/pkg/neo4j"
)

type RAGService struct{}

func (s *RAGService) GetRelevantContent(ctx context.Context, symptoms []string) (string, error) {
	query := `
        MATCH (s:Symptom)
        WHERE s.name IN $symptoms
        MATCH (s)-[:INDICATES]->(c:Condition)
        MATCH (c)-[:LINKED_TO]->(e:EducationalContent)
        RETURN e.content AS content
    `
	params := map[string]interface{}{
		"symptoms": symptoms,
	}

	result, err := neo4j.ExecuteQuery("my-neo4j", query, params)
	if err != nil {
		return "", fmt.Errorf("failed to execute query: %w", err)
	}

	if len(result.Records) == 0 {
		return "No relevant content found.", nil
	}

	content := ""
	for _, record := range result.Records {
		// Capture both return values
		val, ok := record.Get("content")
		if !ok {
			return "", fmt.Errorf("content field not found in record")
		}
		// val is a string if your library returns (string, bool)
		content += val + "\n"
	}

	return content, nil
}
