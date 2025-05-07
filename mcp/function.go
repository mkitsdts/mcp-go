package mcp

import (
	"encoding/json"
	"fmt"
)

// 解析大模型响应信息并调用工具
func (s *MCPClient) parseresp(respBody *[]byte) (string, error) {
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
	if resp.Choices[0].FinishReason == "stop" {
		return resp.Choices[0].Message.Content, nil
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
	s.context = append(s.context, req_mess{Role: "assistant", Content: resp.Choices[0].Message.Content})
	content, err := s.use_tool(resp.Choices[0].Message.Tool_calls[0].Function.Name, args)
	if err != nil {
		fmt.Println("调用工具错误:", err)
		return "", err
	}
	reqBody, err := s.create_request("工具调用结果： "+content, "user")
	if err != nil {
		fmt.Println("转换请求体为JSON错误:", err)
		return "", err
	}
	newRespBody, err := s.send_request(reqBody)
	if err != nil {
		fmt.Println("读取响应体错误:", err)
		return "", err
	}
	return s.parseresp(newRespBody)
}

func (s *MCPClient) use_tool(name string, args map[string]any) (string, error) {
	for i := range *s.golbaltool {
		if (*s.golbaltool)[i].Function.Name == name {
			// 调用工具
			result, err := (*s.golbaltool)[i].Handler(args)
			if err != nil {
				return "", err
			}
			// 返回结果给大模型
			return fmt.Sprintf("%v", result), nil
		}
	}
	for i := range s.tools {
		if s.tools[i].Function.Name == name {
			// 调用工具
			result, err := s.tools[i].Handler(args)
			if err != nil {
				return "", err
			}
			// 返回结果给大模型
			return fmt.Sprintf("%v", result), nil
		}
	}
	return "", fmt.Errorf("error: no tool found with name %s", name)
}
