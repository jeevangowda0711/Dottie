// models.go
// This file contains shared structs for nodes and relationships.

package db

type Node struct {
	Name       string                 `json:"name"`
	Labels     []string               `json:"labels"`
	Properties map[string]interface{} `json:"properties"`
}

type Relationship struct {
	StartNode string `json:"startNode"`
	EndNode   string `json:"endNode"`
	Type      string `json:"type"`
}

type NormalRange struct {
	Name string `json:"name"`
	Min  int    `json:"min"`
	Max  int    `json:"max"`
	Unit string `json:"unit"`
}

type Condition struct {
	Name              string `json:"name"`
	Definition        string `json:"definition"`
	Severity          string `json:"severity"`
	RequiresAttention bool   `json:"requiresAttention"`
}

type Symptom struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

type EducationalContent struct {
	Type   string `json:"type"`
	Url    string `json:"url"`
	Title  string `json:"title"`
	Source string `json:"source"`
}

type ACOGData struct {
	NormalRanges       []NormalRange       `json:"normalRanges"`
	Conditions         []Condition         `json:"conditions"`
	Symptoms           []Symptom           `json:"symptoms"`
	EducationalContent []EducationalContent `json:"educationalContent"`
}
