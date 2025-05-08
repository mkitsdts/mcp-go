package mcp

type Paramaters struct {
	Type       string         `json:"type"`
	Properties map[string]any `json:"properties"`
	Required   []string       `json:"required"`
}

type Tool struct {
	Type     string                                    `json:"type"`
	Function req_fn                                    `json:"function"`
	Handler  func(args map[string]any) (string, error) `json:"-"`
}

type req_fn struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Para        Paramaters `json:"parameters"`
}

type resp_fn struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type tool_call struct {
	Id       string  `json:"id"`
	Tp       string  `json:"type"`
	Function resp_fn `json:"function"`
}
