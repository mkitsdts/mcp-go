package mcp

import (
	"fmt"
	"net/http"
	"os"
	"slices"
)

type MCPClient struct {
	index   int
	client  []*http.Client
	context []req_mess
	files   map[string]string
	tools   []Tool
	host    string
	key     string
	name    string
}

func (s *MCPClient) GetHTTPClient() *http.Client {
	s.index = (s.index + 1) % len(s.client)
	return s.client[s.index]
}

func NewMCPClient(name string, host string, key string) *MCPClient {
	return &MCPClient{
		client:  make([]*http.Client, MAX_HTTP_CLIENT_CONNECTIONS),
		context: make([]req_mess, 0),
		files:   make(map[string]string),
		tools:   make([]Tool, 0),
		host:    host,
		key:     key,
		name:    name,
	}
}

const MAX_FILE_SIZE int64 = 10 * 1024 * 1024 // 10MB

func (s *MCPClient) AddFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("invalid file path")
	}
	fileSize := fileInfo.Size()
	if fileSize > MAX_FILE_SIZE {
		return fmt.Errorf("exceeded maximum file size (10MB)")
	}
	// 读取文件内容
	// 转化为字符串
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	// 将文件内容转换为字符串
	content := string(fileContent)
	s.files[fileInfo.Name()] = content
	return nil
}

func (s *MCPClient) EraseFile(name string) {
	for i := range s.files {
		if i == name {
			delete(s.files, i)
			break
		}
	}
}

func (s *MCPClient) ClearFile() {
	s.files = make(map[string]string)
}

func (s *MCPClient) Chat(context string) (string, error) {
	// 提取信息
	extractKeywordBodyJSON, err := s.create_request(context, "user")
	if err != nil {
		return "", err
	}
	// 发送 POST 请求
	keywordBody, err := s.send_request(extractKeywordBodyJSON)
	if err != nil {
		return "", err
	}
	// 解析结果并调用工具
	result, err := s.parseresp(keywordBody)
	if err != nil {
		if err.Error() == "error: no tool calls found in response" {
			return result, nil
		}
		return "", err
	}
	return result, nil
}

func (s *MCPClient) AddTool(name string, description string, parameters Paramaters, handler func(args map[string]any) (string, error)) error {
	if len(s.tools) > 10 {
		return fmt.Errorf("exceeded maximum number of tools (10)")
	}
	for i := range s.tools {
		if s.tools[i].Function.Name == name {
			return nil
		}
	}
	tool := Tool{
		Type: "function",
		Function: struct {
			Name        string     `json:"name"`
			Description string     `json:"description"`
			Para        Paramaters `json:"parameters"`
		}{
			Name:        name,
			Description: description,
			Para:        parameters,
		},
		Handler: handler,
	}
	s.tools = append(s.tools, tool)
	return nil
}

func (s *MCPClient) EraseTool(name string) {
	for i := range s.tools {
		if s.tools[i].Function.Name == name {
			s.tools = slices.Delete(s.tools, i, i+1)
			break
		}
	}
}

const MAX_CONTEXT_LENGTH = 4096

func (s *MCPClient) AddSystemPrompt(context string) error {
	if len(s.context) > MAX_CONTEXT_LENGTH {
		return fmt.Errorf("exceeded maximum context length (%d)", MAX_CONTEXT_LENGTH)
	}
	s.context = append(s.context, req_mess{Role: "system", Content: context})
	return nil
}

func (s *MCPClient) ClearTool() {
	s.tools = s.tools[:0]
}

func (s *MCPClient) ClearHistory() {
	s.context = s.context[:0]
}
