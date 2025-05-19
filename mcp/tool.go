package mcp

import (
	"encoding/json"
	"fmt"
)

// 解析大模型响应信息并调用工具
func (s *MCPClient) parseresp(respBody *[]byte) (string, error) {
	max_times := 10
	curr_times := 0

	for {
		if curr_times >= max_times {
			return "", fmt.Errorf("error: maximum retry limit reached")
		}
		curr_times++
		resp := response{}
		if err := json.Unmarshal(*respBody, &resp); err != nil {
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
		fmt.Println("finish reason:", resp.Choices[0].FinishReason)

		// 检查是否有工具调用
		if len(resp.Choices[0].Message.Tool_calls) == 0 {
			if resp.Choices[0].Message.Content != "" {
				s.context = append(s.context, req_mess{Role: "assistant", Content: resp.Choices[0].Message.Content})
			}
			return "", nil
		}
		// 解析参数
		var args map[string]any
		if resp.Choices[0].Message.Tool_calls[0].Function.Name == "" {
			return "", fmt.Errorf("error: no function name found in tool call")
		}
		if resp.Choices[0].Message.Tool_calls[0].Function.Arguments == "" {
			return "", fmt.Errorf("error: no function arguments found in tool call")
		}
		args = make(map[string]any)
		if err := json.Unmarshal([]byte(resp.Choices[0].Message.Tool_calls[0].Function.Arguments), &args); err != nil {
			return "", err
		}
		result, err := s.use_tool(resp.Choices[0].Message.Tool_calls[0].Function.Name, args)
		if err != nil {
			return "", err
		}
		s.context = append(s.context, req_mess{Role: "user", Content: tool_response_prompt + result})
		reqBody, err := s.create_request("工具调用结果如上，请总结答案或再次调用工具", "user")
		if err != nil {
			return "", err
		}
		respBody, err = s.send_request(reqBody)
		fmt.Println("respBody:", string(*respBody))
		if err != nil {
			return "", err
		}
	}

}

func (s *MCPClient) use_tool(name string, args map[string]any) (string, error) {
	// 检查工具名称是否存在
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
