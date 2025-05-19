package mcp

import "time"

const (
	MAX_CLIENT_CONNECTION_TIME = 180 * time.Second
	MAX_CLIENT_IDLE_TIME       = 60 * time.Second
	MAX_CLIENT_IDLE_CONNS      = 10
)

const (
	MAX_HTTP_CLIENT_CONNECTIONS = 50
)
