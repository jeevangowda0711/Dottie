// schema.go
// This file initializes the Neo4j schema and populates it with nodes and relationships.

package db

import (
	"fmt"

	"github.com/hypermodeinc/modus/sdk/go/pkg/neo4j"
)

// CreateNodes initializes the Neo4j database with predefined nodes
func CreateNodes() error {
	// Define normal parameters
	normalRanges := []Node{
		{
			Name:   "CycleLength",
			Labels: []string{"Parameter"},
			Properties: map[string]interface{}{
				"min":  21,
				"max":  45,
				"unit": "days",
			},
		},
		{
			Name:   "CycleDuration",
			Labels: []string{"Parameter"},
			Properties: map[string]interface{}{
				"min":  3,
				"max":  7,
				"unit": "days",
			},
		},
	}

	// Define conditions
	conditions := []Node{
		{
			Name:   "Amenorrhea",
			Labels: []string{"Condition"},
			Properties: map[string]interface{}{
				"definition":        "No periods for >3 months",
				"severity":          "high",
				"requiresAttention": true,
			},
		},
		{
			Name:   "Oligomenorrhea",
			Labels: []string{"Condition"},
			Properties: map[string]interface{}{
				"definition":        "Cycles >45 days apart",
				"severity":          "moderate",
				"requiresAttention": true,
			},
		},
	}

	// Combine all nodes
	allNodes := append(normalRanges, conditions...)

	// Create nodes in Neo4j
	for _, node := range allNodes {
		query := `
		CREATE (n:` + node.Labels[0] + ` {name: $name, properties: $properties})
		`
		params := map[string]interface{}{
			"name":       node.Name,
			"properties": node.Properties,
		}
		if _, err := neo4j.ExecuteQuery("my-neo4j", query, params); err != nil {
			return fmt.Errorf("failed to create node: %w", err)
		}
		fmt.Printf("Node created: %s\n", node.Name)
	}

	fmt.Println("Schema initialized successfully!")
	return nil
}
