package mcp

type resp_message struct {
	Role       string      `json:"role"`
	Tool_calls []tool_call `json:"tool_calls"`
	Content    string      `json:"content"`
}

type choice struct {
	Index        int          `json:"index"`
	Logprobs     string       `json:"logprobs"`
	FinishReason string       `json:"finish_reason"`
	Message      resp_message `json:"message"`
}

type response struct {
	Choices     []choice `json:"choices"`
	Temperature float32  `json:"temperature"`
	Error       string   `json:"error"`
}

type req_mess struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type request struct {
	Model       string  `json:"model"`
	Messages    any     `json:"messages"`
	Temperature float32 `json:"temperature"`
	Stream      bool    `json:"stream"`
	Tools       []Tool  `json:"tools"`
}
