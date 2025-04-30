package mcp

import "net/http"

// 核心服务
// 包含上下文协议、工具、消息等
type MCPService struct {
	clients []MCPClient
	host    string
	name    string
	key     string
}

type MCPClient struct {
	client  http.Client
	context []req_mess
	tools   []Tool
	host    *string
	key     *string
	name    *string
}
