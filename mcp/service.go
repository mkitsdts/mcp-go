package mcp

import "net/http"

// 核心服务
// 包含上下文协议、工具、消息等
type MCPService struct {
	Clients map[string]MCPClient
	tools   []Tool
	host    string
	name    string
	key     string
}

type MCPClient struct {
	client     http.Client
	context    []req_mess
	tools      []Tool
	host       *string
	key        *string
	name       *string
	golbaltool *[]Tool
}
