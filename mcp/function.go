package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (s *MCPService) extract_keyword(context string) ([]byte, error) {
	messages := []map[string]string{
		{"role": "system", "content": system_extract_paramater_prompt},
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
	fmt.Println("Request body:", string(requestBodyJSON))
	return requestBodyJSON, nil
}

func (s *MCPService) get_answer(respBody *[]byte) (string, error) {
	// 解析 JSON 响应
	var result map[string]any
	if err := json.Unmarshal(*respBody, &result); err != nil {
		fmt.Println("解析 JSON 响应错误:", err)
		return "", err
	}
	// 打印响应内容
	fmt.Println("响应内容:", string(*respBody))
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
	choices, ok := result["choices"]
	if !ok {
		return "", fmt.Errorf("error: no choices found in response")
	}
	choicesArray, ok := choices.([]any)
	if !ok || len(choicesArray) == 0 {
		return "", fmt.Errorf("error: no choices found in response")
	}
	choice, ok := choicesArray[0].(map[string]any)
	if !ok {
		return "", fmt.Errorf("error: invalid choice format in response")
	}
	finishReason, ok := choice["finish_reason"]
	if !ok || finishReason != "tool_calls" {
		return "", fmt.Errorf("error: finish reason is not tool_calls")
	}
	message, ok := choice["message"]
	if !ok {
		return "", fmt.Errorf("error: no message found in choice")
	}
	messageMap, ok := message.(map[string]any)
	if !ok {
		return "", fmt.Errorf("error: invalid message format in choice")
	}
	toolCall, ok := messageMap["tool_calls"]
	if !ok {
		return "", fmt.Errorf("error: no tool call found in message")
	}
	toolCallArray, ok := toolCall.([]any)
	if !ok || len(toolCallArray) == 0 {
		return "", fmt.Errorf("error: no tool call found in message")
	}
	toolCallMap, ok := toolCallArray[0].(map[string]any)
	if !ok {
		return "", fmt.Errorf("error: invalid tool call format in message")
	}
	function, ok := toolCallMap["function"]
	if !ok {
		return "", fmt.Errorf("error: no function found in message")
	}
	functionName, ok := function.(map[string]any)["name"]
	if !ok {
		return "", fmt.Errorf("error: no function name found in message")
	}
	functionArguments, ok := function.(map[string]any)["arguments"]
	if !ok {
		return "", fmt.Errorf("error: no function arguments found in message")
	}
	// 调用工具
	for i := range s.Tools {
		if s.Tools[i].Function.Name == functionName {
			// 解析参数
			var args map[string]any
			if err := json.Unmarshal([]byte(functionArguments.(string)), &args); err != nil {
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

	return "", nil
}

func (s *MCPService) extract_result(result string) (string, error) {
	// 把工具调用结果发送给大模型
	fmt.Println("工具最初返回结果:", result)
	messages := []map[string]string{
		{"role": "system", "content": system_extarct_answer_prompt},
		{"role": "user", "content": s.Context},
		{"role": "tool", "content": "工具结果:" + result},
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
	}
	// 转换为JSON
	requestBodyJSON, err := json.Marshal(requestBody)
	fmt.Println("Request body:", string(requestBodyJSON))
	if err != nil {
		fmt.Println("转换请求体为JSON错误:", err)
		return "", err
	}
	// 发送请求
	resp, err := s.Client.Post(s.Host, request_content_type, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		fmt.Println("发送请求错误:", err)
		return "", err
	}
	defer resp.Body.Close()
	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应体错误:", err)
		return "", err
	}
	// 转换成字符串
	resultString := string(respBody)

	// 解析 JSON 响应
	var resultMap map[string]any
	if err := json.Unmarshal([]byte(resultString), &resultMap); err != nil {
		fmt.Println("解析 JSON 响应错误:", err)
		return "", err
	}
	// 检查是否有响应内容
	errMsg, ok := resultMap["error"]
	if ok {
		if errMsgMap, ok := errMsg.(map[string]any); ok {
			if message, ok := errMsgMap["message"]; ok {
				fmt.Println("返回错误消息:", message)
				return "", fmt.Errorf("error: %s", message)
			}
		}
		fmt.Println("返回错误消息:", errMsg)
		return "", fmt.Errorf("error: %s", errMsg)
	}
	type Result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	var resultObj Result
	if err := json.Unmarshal([]byte(resultString), &resultObj); err != nil {
		fmt.Println("解析 JSON 响应错误:", err)
		return "", err
	}
	if len(resultObj.Choices) > 0 {
		return resultObj.Choices[0].Message.Content, nil
	}
	fmt.Println("没有找到有效的响应内容")
	return "", nil
}

func (s *MCPService) try_get_model_info() bool {
	resp, err := s.Client.Get(s.Host + s.Name)
	if err != nil {
		return false
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return false
	}
	fmt.Println(string(body))
	return resp.StatusCode == http.StatusOK
}

func (s *MCPService) sendContextToModel(data *[]byte) (*[]byte, error) {
	// 发送 POST 请求
	req, err := http.NewRequest("POST", s.Host, bytes.NewBuffer(*data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", request_content_type)
	if s.Key != "" {
		req.Header.Set("Authorization", "Bearer "+s.Key)
	}
	resp, err := s.Client.Do(req)
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
