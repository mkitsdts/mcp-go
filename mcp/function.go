package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// 提取回答关键信息
func (s *MCPClient) create_extract_keyword_request(context string) (*[]byte, error) {
	s.context = append(s.context, req_mess{Role: "user", Content: context})

	// 转换为JSON
	messagesJSON, err := json.Marshal(s.context)
	if err != nil {
		return nil, err
	}

	// 创建完整请求体
	requestBody := request{
		Model:       *s.name,
		Messages:    json.RawMessage(messagesJSON),
		Temperature: 0.1,
		Stream:      false,
		Tools:       s.tools,
	}
	// 转换为JSON
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("转换请求体为JSON错误:", err)
		return nil, err
	}
	fmt.Println("Request body:", string(requestBodyJSON))
	return &requestBodyJSON, nil
}

// 解析大模型响应信息并调用工具
func (s *MCPClient) get_tool(respBody *[]byte) (string, error) {
	resp := response{}
	if err := json.Unmarshal(*respBody, &resp); err != nil {
		fmt.Println("解析 JSON 响应错误:", err)
		return "", err
	}
	// 检查是否有响应内容
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("error: no choices found in response")
	}
	// 检查是否有错误
	if resp.Error != "" {
		return "", fmt.Errorf("error: %s", resp.Error)
	}
	// 检查是否有工具调用
	if len(resp.Choices[0].Message.Tool_calls) == 0 {
		return "", fmt.Errorf("error: no tool calls found in response")
	}
	// 解析参数
	var args map[string]any
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Tool_calls[0].Function.Arguments), &args); err != nil {
		fmt.Println("转换函数参数为JSON错误:", err)
		return "", err
	}
	// 调用工具
	for i := range s.tools {
		if s.tools[i].Function.Name == resp.Choices[0].Message.Tool_calls[0].Function.Name {
			// 调用工具
			result, err := s.tools[i].Handler(args)
			if err != nil {
				return "", err
			}
			// 返回结果给大模型
			return fmt.Sprintf("%v", result), nil
		}
	}

	return "", nil
}

// 把工具调用结果发送给大模型
func (s *MCPClient) create_extract_result_request(result string) (*[]byte, error) {
	s.context = append(s.context, req_mess{Role: "user", Content: "工具结果： " + result})
	// 转换为JSON
	messagesJSON, err := json.Marshal(s.context)
	if err != nil {
		return nil, err
	}
	// 创建完整请求体
	requestBody := request{
		Model:       *s.name,
		Messages:    json.RawMessage(messagesJSON),
		Temperature: 0.1,
		Stream:      false,
	}
	// 转换为JSON
	requestBodyJSON, err := json.Marshal(requestBody)
	fmt.Println("Request body:", string(requestBodyJSON))
	if err != nil {
		fmt.Println("转换请求体为JSON错误:", err)
		return nil, err
	}
	return &requestBodyJSON, nil
}

func (s *MCPClient) get_result(respBody *[]byte) (string, error) {
	// 解析 JSON 响应
	resp := response{}
	if err := json.Unmarshal(*respBody, &resp); err != nil {
		fmt.Println("解析 JSON 响应错误:", err)
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

const request_content_type string = "application/json"

// 向大模型发送请求
func (s *MCPClient) send_request(data *[]byte) (*[]byte, error) {
	// 发送 POST 请求
	req, err := http.NewRequest("POST", (*s.host), bytes.NewBuffer(*data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", request_content_type)
	if (*s.key) != "" {
		req.Header.Set("Authorization", "Bearer "+(*s.key))
	}
	fmt.Println("请求头:", req.Header)
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应体错误:", err)
		return nil, err
	}
	return &respBody, nil
}
