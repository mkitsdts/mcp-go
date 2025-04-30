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
	s := mcp.NewMCPService("deepseek-chat", "https://api.deepseek.com/chat/completions", mcp.API_KEY)
	dialog := s.NewDialogue()
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

	// 测试1：直接查询天气
	fmt.Println("\n--- 测试1：查询北京天气 ---")
	resp1, err := dialog.Chat("北京今天的天气怎么样？")
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}
	fmt.Printf("模型回复: %s\n", resp1)

	fmt.Println("\n所有测试完成")
}
