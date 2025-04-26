package mcp

import (
	"fmt"
	"io"
	"net/http"
)

func (s *MCPService) try_get_model_info() bool {
	resp, err := s.Client.Get(s.Host + "/v1/models/" + s.Name)
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
