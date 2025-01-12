package graph

import (
	"context"
	"fmt"

	"github.com/hypermodeinc/modus/sdk/go/pkg/neo4j"
)

type Content struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Category string `json:"category"`
	Title    string `json:"title"`
	Source   string `json:"source"`
	URL      string `json:"url"`
	Abstract string `json:"abstract"`
}

func GetEducationalContent(ctx context.Context, condition string) ([]Content, error) {
	query := `
    MATCH (c:Condition {name: $condition})-[:LINKED_TO]->(e:EducationalContent)
    RETURN e.id AS id, e.type AS type, e.category AS category, e.title AS title, e.source AS source, e.url AS url, e.abstract AS abstract
    `
	params := map[string]interface{}{
		"condition": condition,
	}

	result, err := neo4j.ExecuteQuery("my-neo4j", query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch educational content: %w", err)
	}

	var contents []Content
	for _, record := range result.Records {
		// 1) Get the string, 2) Check if it exists
		id, ok := record.Get("id")
		if !ok {
			return nil, fmt.Errorf("id field not found in record")
		}

		contentType, ok := record.Get("type")
		if !ok {
			return nil, fmt.Errorf("type field not found in record")
		}

		category, ok := record.Get("category")
		if !ok {
			return nil, fmt.Errorf("category field not found in record")
		}

		title, ok := record.Get("title")
		if !ok {
			return nil, fmt.Errorf("title field not found in record")
		}

		source, ok := record.Get("source")
		if !ok {
			return nil, fmt.Errorf("source field not found in record")
		}

		url, ok := record.Get("url")
		if !ok {
			return nil, fmt.Errorf("url field not found in record")
		}

		abstract, ok := record.Get("abstract")
		if !ok {
			return nil, fmt.Errorf("abstract field not found in record")
		}

		// Now that we have them, each is a string (assuming your Get returns (string, bool))
		content := Content{
			ID:       id,
			Type:     contentType,
			Category: category,
			Title:    title,
			Source:   source,
			URL:      url,
			Abstract: abstract,
		}
		contents = append(contents, content)
	}

	return contents, nil
}
