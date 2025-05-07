package mcp

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func (s *MCPClient) Chat(context string) (string, error) {
	// 提取信息
	extractKeywordBodyJSON, err := s.create_request(context, "user")
	if err != nil {
		fmt.Println("转换请求体为JSON错误:", err)
		return "", err
	}
	// 发送 POST 请求
	keywordBody, err := s.send_request(extractKeywordBodyJSON)
	if err != nil {
		fmt.Println("读取响应体错误:", err)
		return "", err
	}
	fmt.Println("响应内容:", string(*keywordBody))
	// 解析结果并调用工具
	result, err := s.parseresp(keywordBody)
	if err != nil {
		if err.Error() == "error: no tool calls found in response" {
			fmt.Println("没有工具调用，直接返回响应内容")
			return result, nil
		}
		fmt.Println("解析响应结果错误:", err)
		return "", err
	}
	return result, nil
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

const MAX_FILE_SIZE int64 = 10 * 1024 * 1024 // 10MB

func (s *MCPClient) AddFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("打开文件错误:", err)
		return err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("获取文件信息错误:", err)
		return err
	}
	if fileInfo.IsDir() {
		fmt.Println("指定的路径是一个目录，无法添加文件")
		return fmt.Errorf("指定的路径是一个目录，无法添加文件")
	}
	fileSize := fileInfo.Size()
	if fileSize > MAX_FILE_SIZE {
		fmt.Println("文件大小超过限制，无法添加文件")
		return fmt.Errorf("文件大小超过限制，无法添加文件")
	}
	// 读取文件内容
	// 转化为字符串
	fileContent, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("读取文件内容错误:", err)
		return err
	}
	// 将文件内容转换为字符串
	content := string(fileContent)
	fmt.Println("文件内容:", content)
	s.files[fileInfo.Name()] = content
	return nil
}

func (s *MCPService) AddGolbalFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("打开文件错误:", err)
		return err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("获取文件信息错误:", err)
		return err
	}
	if fileInfo.IsDir() {
		fmt.Println("指定的路径是一个目录，无法添加文件")
		return fmt.Errorf("指定的路径是一个目录，无法添加文件")
	}
	fileSize := fileInfo.Size()
	if fileSize > MAX_FILE_SIZE {
		fmt.Println("文件大小超过限制，无法添加文件")
		return fmt.Errorf("文件大小超过限制，无法添加文件")
	}
	// 读取文件内容
	// 转化为字符串
	fileContent, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("读取文件内容错误:", err)
		return err
	}
	// 将文件内容转换为字符串
	content := string(fileContent)
	fmt.Println("文件内容:", content)
	s.files[fileInfo.Name()] = content
	return nil
}

func (s *MCPClient) AddTool(name string, description string, parameters Paramaters, handler func(args map[string]any) (string, error)) {
	if len(s.tools)+len(*s.golbaltool) > 10 {
		fmt.Println("工具数量超过限制，无法添加新工具")
		return
	}
	for i := range *s.golbaltool {
		if (*s.golbaltool)[i].Function.Name == name {
			fmt.Println("工具名称已存在，无法添加新工具")
			return
		}
	}
	for i := range s.tools {
		if s.tools[i].Function.Name == name {
			fmt.Println("工具名称已存在，无法添加新工具")
			return
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
	fmt.Println("Tool", tool)
}

func (s *MCPClient) ClearHistory() {
	s.context = s.context[:0]
}

func (s *MCPService) NewClient(tag string) *MCPClient {
	c := MCPClient{}
	c.client = http.Client{}
	c.client.Timeout = 180 * time.Second
	c.client.Transport = &http.Transport{
		MaxIdleConns:          10,
		IdleConnTimeout:       60 * time.Second,
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

func (s *MCPService) GetClient(tag string) *MCPClient {
	c, ok := s.clients[tag]
	if !ok {
		fmt.Println("没有找到指定的客户端")
		return nil
	}
	return &c
}
