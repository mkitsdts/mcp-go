package lmstudio

import "net/http"

type Para struct {
	Type       string         `json:"type"`
	Properties map[string]any `json:"properties"`
	Required   []string       `json:"required"`
}

type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  Para   `json:"parameters"`
	Handler     func(args ...any) (any, error)
}

type MCPService struct {
	Client http.Client
	Host   string
	Name   string
	Key    string
	Tools  []Tool
}
