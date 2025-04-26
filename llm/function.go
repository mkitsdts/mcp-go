package mcp

import (
	"encoding/json"
	"fmt"
)

func (s *MCPService) extract_keyword(context string) ([]byte, error) {
	messages := []map[string]string{
		{"role": "system", "content": system_prompt},
		{"role": "user", "content": context},
	}
	fmt.Println("Request data:", messages)
	// 转换为JSON
	messagesJSON, err := json.Marshal(messages)
	if err != nil {
		return nil, err
	}

	// 创建完整请求体
	requestBody := map[string]any{
		"model":       s.Name,
		"messages":    json.RawMessage(messagesJSON),
		"temperature": 0.1,
		"stream":      false,
		"tools":       s.Tools,
	}
	// 转换为JSON
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("转换请求体为JSON错误:", err)
		return nil, err
	}
	return requestBodyJSON, nil
}

func (s *MCPService) extract_result(respBody []byte) (string, error) {
	// 解析 JSON 响应
	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		fmt.Println("解析 JSON 响应错误:", err)
		return "", err
	}
	// 打印响应内容
	fmt.Println("响应内容:", string(respBody))
	// 检查是否有响应内容
	if errMsg, ok := result["error"]; ok {
		if errMsgMap, ok := errMsg.(map[string]any); ok {
			if message, ok := errMsgMap["message"]; ok {
				fmt.Println("返回错误消息:", message)
				return "", fmt.Errorf("error: %s", message)
			}
		}
	}
	// 检查是否调用工具
	if choices, ok := result["choices"]; ok {
		if choicesArray, ok := choices.([]any); ok && len(choicesArray) > 0 {
			if choice, ok := choicesArray[0].(map[string]any); ok {
				if finishReason, ok := choice["finish_reason"]; ok && finishReason == "tool_calls" {
					if message, ok := choice["message"]; ok {
						if messageMap, ok := message.(map[string]any); ok {
							if functionCall, ok := messageMap["function_call"]; ok {
								if functionCallMap, ok := functionCall.(map[string]any); ok {
									if name, ok := functionCallMap["name"]; ok {
										if arguments, ok := functionCallMap["arguments"]; ok {
											// 调用工具
											for i := range s.Tools {
												if s.Tools[i].Function.Name == name {
													// 解析参数
													var args map[string]any
													if err := json.Unmarshal([]byte(arguments.(string)), &args); err != nil {
														fmt.Println("解析参数错误:", err)
														return "", err
													}
													// 调用工具
													result, err := s.Tools[i].Handler(args)
													if err != nil {
														fmt.Println("调用工具错误:", err)
														return "", err
													}
													// 返回结果给大模型

													return fmt.Sprintf("%v", result), nil
												}
											}
										}
									}
								}
							}
						}
					}
				} else if content, ok := choice["content"]; ok {
					return content.(string), nil
				}
			}
		}
	}
	return "", nil
}
