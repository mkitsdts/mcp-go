package mcp

import (
	"fmt"
	"net/http"
	"os"
	"slices"
	"time"
)

// 核心服务
// 包含上下文协议、工具、消息等
type MCPService struct {
	clients map[string]MCPClient
	tools   []Tool
	files   map[string]string
	host    string
	name    string
	key     string
}

func NewMCPService(name string, host string, key string) *MCPService {
	s := &MCPService{}
	s.name = name
	s.host = host
	s.key = key
	s.clients = map[string]MCPClient{}
	s.tools = []Tool{}
	s.files = map[string]string{}
	s.clients = make(map[string]MCPClient)
	s.tools = make([]Tool, 0)
	return s
}

func (s *MCPService) AddGlobalTool(name string, description string, parameters Paramaters, handler func(args map[string]any) (string, error)) {
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
}

func (s *MCPService) EraseGlobalTool(name string) {
	for i := range s.tools {
		if s.tools[i].Function.Name == name {
			s.tools = slices.Delete(s.tools, i, i+1)
			break
		}
	}
}

func (s *MCPService) ClearGlobalTool() {
	s.tools = s.tools[:0]
}

func (s *MCPService) AddGolbalFile(path string) error {
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
		return fmt.Errorf("file too large")
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

func (s *MCPService) EraseGolbalFile(name string) {
	for i := range s.files {
		if i == name {
			delete(s.files, i)
			break
		}
	}
}

func (s *MCPService) ClearGolbalFile() {
	s.files = make(map[string]string)
}

func (s *MCPService) NewClient(tag string) *MCPClient {
	c := MCPClient{}
	c.client = http.Client{}
	c.client.Timeout = MAX_CLIENT_CONNECTION_TIME
	c.client.Transport = &http.Transport{
		MaxIdleConns:          MAX_CLIENT_IDLE_CONNS,
		IdleConnTimeout:       MAX_CLIENT_IDLE_TIME,
		DisableCompression:    true,
		ForceAttemptHTTP2:     true,
		ExpectContinueTimeout: 10 * time.Second,
	}
	c.context = []req_mess{}
	c.context = append(c.context, req_mess{Role: "system", Content: system_prompt})
	c.tools = []Tool{}
	c.host = &s.host
	c.key = &s.key
	c.name = &s.name
	c.golbaltool = &s.tools
	c.golbalfile = &s.files
	c.files = map[string]string{}
	s.clients[tag] = c
	return &c
}

func (s *MCPService) EraseClient(tag string) {
	c, ok := s.clients[tag]
	if !ok {
		return
	}
	c.context = []req_mess{}
	c.files = map[string]string{}
	c.tools = []Tool{}
	c.host = nil
	c.key = nil
	c.name = nil
	c.golbaltool = nil
	c.golbalfile = nil
	c.client.CloseIdleConnections()
	c.client.Transport = nil
	c.client = http.Client{}
	c.client.Timeout = MAX_CLIENT_CONNECTION_TIME
	delete(s.clients, tag)
}

func (s *MCPService) GetClient(tag string) *MCPClient {
	c, ok := s.clients[tag]
	if !ok {
		return nil
	}
	return &c
}
