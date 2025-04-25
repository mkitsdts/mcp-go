package lmstudio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func (s *MCPService) Chat(context string) (string, error) {
	// 把字符串转换成io.Reader

	messages := []map[string]string{
		{"role": "system", "content": system_prompt},
		{"role": "user", "content": context},
	}
	fmt.Println("Request data:", messages)
	// 转换为JSON
	messagesJSON, err := json.Marshal(messages)
	if err != nil {
		return "", err
	}

	// 创建完整请求体
	requestBody := map[string]any{
		"model":       s.Name,
		"messages":    json.RawMessage(messagesJSON),
		"temperature": 0.1,
		"stream":      false,
		"tools": []map[string]any{
			{
				"type": "function", // 必须添加这个字段，表明这是函数类型的工具
				"function": map[string]any{
					"name":        "weather_query",
					"description": "search weather by location",
					"parameters": map[string]any{ // 不是inputSchema，应该是parameters
						"type": "object", // 需要指定这个字段
						"properties": map[string]any{
							"latitude": map[string]any{
								"type":        "number",
								"description": "latitude coordinate",
							},
							"longitude": map[string]any{
								"type":        "number",
								"description": "longitude coordinate",
							},
						},
						"required": []string{"latitude", "longitude"},
					},
				},
			},
		},
	}
	// 转换为JSON
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("转换请求体为JSON错误:", err)
		return "", err
	}

	// 发送 POST 请求
	resp, err := s.Client.Post(s.Host+"/v1/chat/completions", "application/json", bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应体错误:", err)
		return "", err
	}
	// 解析 JSON 响应
	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		fmt.Println("解析 JSON 响应错误:", err)
		return "", err
	}
	// 打印响应内容
	fmt.Println("响应内容:", string(respBody))
	// 检查是否有错误信息
	if errMsg, ok := result["error"]; ok {
		if errMsgMap, ok := errMsg.(map[string]any); ok {
			if message, ok := errMsgMap["message"]; ok {
				fmt.Println("返回错误消息:", message)
				return "", fmt.Errorf("error: %s", message)
			}
		}
	}
	// 检查是否有响应内容
	if choices, ok := result["choices"]; ok {
		if choicesArray, ok := choices.([]any); ok && len(choicesArray) > 0 {
			if choice, ok := choicesArray[0].(map[string]any); ok {
				if message, ok := choice["message"]; ok {
					if toolCalls, ok := message.(map[string]any)["tool_calls"]; ok {
						if function, ok := toolCalls.([]any)[0].(map[string]any)["function"]; ok {
							if arguments, ok := function.(map[string]any)["arguments"]; ok {
								// 解析函数参数
								argumentsJSON, err := json.Marshal(arguments)
								if err != nil {
									fmt.Println("解析函数参数错误:", err)
									return "", err
								}
								// 将函数参数转换为字符串
								argumentsStr, err := json.MarshalIndent(arguments, "", "  ")
								if err != nil {
									fmt.Println("转换函数参数为字符串错误:", err)
									return "", err
								}
								// 打印函数参数
								fmt.Println("函数参数:", string(argumentsStr))
								// 解析函数参数为 map
								var functionArgs map[string]any
								if err := json.Unmarshal(argumentsJSON, &functionArgs); err != nil {
									fmt.Println("解析函数参数为 map 错误:", err)
									return "", err
								}
								// 获取经纬度
								latitude, ok := functionArgs["latitude"].(float64)
								if !ok {
									fmt.Println("获取纬度错误")
									return "", fmt.Errorf("error: 获取纬度错误")
								}
								longitude, ok := functionArgs["longitude"].(float64)
								if !ok {
									fmt.Println("获取经度错误")
									return "", fmt.Errorf("error: 获取经度错误")
								}
								// 打印经纬度
								fmt.Printf("经度: %f, 纬度: %f\n", longitude, latitude)
							}
						} else {
							fmt.Println("没有找到函数调用")
						}
					}
				}
			}
		}
	}
	// 如果没有找到响应内容，返回一个空字符串
	return "", nil
}
