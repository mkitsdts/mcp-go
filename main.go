package main

import (
	"fmt"
	"log"
	mcp "mcp-go/mcp"
)

// 模拟天气查询功能
func getWeather(args map[string]any) (string, error) {
	// 从参数中提取位置信息
	latitude, latOk := args["latitude"].(float64)
	longitude, lonOk := args["longitude"].(float64)

	if !latOk || !lonOk {
		return "", fmt.Errorf("缺少有效的经纬度参数")
	}

	fmt.Println("实际经纬度:", latitude, longitude)

	// 模拟天气数据
	var weather string
	if latitude > 30 && longitude > 100 {
		weather = "晴天，温度28°C，湿度45%"
	} else if latitude > 0 {
		weather = "多云，温度22°C，湿度60%"
	} else {
		weather = "小雨，温度18°C，湿度85%"
	}

	fmt.Printf("获取位置 (%.2f, %.2f) 的天气: %s\n", latitude, longitude, weather)
	return fmt.Sprintf("当前天气: %s", weather), nil
}

func main() {
	// 初始化服务
	s := mcp.NewMCPService("qwen3-14b", "http://localhost:1234/v1/chat/completions", "")
	dialog := s.NewClient("查询天气")
	// 添加天气查询工具
	dialog.AddTool(
		"weather_query",
		"查询指定位置的天气情况",
		mcp.Paramaters{
			Type: "object",
			Properties: map[string]any{
				"latitude": map[string]any{
					"type":        "number",
					"description": "latitude",
				},
				"longitude": map[string]any{
					"type":        "number",
					"description": "longitude",
				},
			},
			Required: []string{"latitude", "longitude"},
		},
		getWeather,
	)
	// 添加华氏度转换工具
	dialog.AddTool(
		"celsius_to_fahrenheit",
		"将摄氏度转换为华氏度",
		mcp.Paramaters{
			Type: "object",
			Properties: map[string]any{
				"celsius": map[string]any{
					"type":        "number",
					"description": "摄氏度",
				},
			},
			Required: []string{"celsius"},
		},
		func(args map[string]any) (string, error) {
			celsius, ok := args["celsius"].(float64)
			if !ok {
				return "", fmt.Errorf("缺少有效的摄氏度参数")
			}
			fahrenheit := celsius*9/5 + 32
			fmt.Printf("摄氏度: %.2f, 华氏度: %.2f\n", celsius, fahrenheit)
			return fmt.Sprintf("华氏度: %.2f", fahrenheit), nil
		},
	)

	// 测试1：直接查询天气
	fmt.Println("\n--- 测试1：查询北京天气 ---")
	resp1, err := dialog.Chat("我要查询北京的天气，然后把摄氏度转换成华氏度")
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}
	fmt.Printf("模型回复: %s\n", resp1)

	dialog.ClearHistory()

	// 测试2:无需工具调用
	fmt.Println("\n--- 测试2：问候 ---")
	resp2, err := dialog.Chat("你好")
	if err != nil {
		log.Fatalf("操作失败: %v", err)
	}
	fmt.Printf("模型回复: %s\n", resp2)

	fmt.Println("\n所有测试完成")
}
