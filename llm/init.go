package lmstudio

import (
	"fmt"
	"io"
	"net/http"
	"time"
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

func (s *MCPService) Init(name string, host string) {
	s.Name = name
	s.Host = host
	s.Client = http.Client{}
	s.Client.Timeout = 60 * time.Second
	s.Client.Transport = &http.Transport{
		MaxIdleConns:          10,
		IdleConnTimeout:       30 * time.Second,
		DisableCompression:    true,
		ForceAttemptHTTP2:     true,
		ExpectContinueTimeout: 10 * time.Second,
	}
	if !s.try_get_model_info() {
		panic("Failed to connect to model server")
	} else {
		fmt.Println("Model server connected successfully")
	}
}
