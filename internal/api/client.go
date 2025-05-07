package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"lazydebrid/internal/config"
	"lazydebrid/internal/logs"
)

var client = &http.Client{Timeout: 10 * time.Second}

func NewRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		logs.LogEvent(fmt.Errorf("failed creating request: %w", err))
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.APIToken()))
	return req, nil
}

func DoRequest(req *http.Request) ([]byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		logs.LogEvent(fmt.Errorf("request error: %w", err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logs.LogEvent(fmt.Errorf("failed reading response: %w", err))
		return nil, err
	}

	if resp.StatusCode >= 400 {
		logs.LogEvent(fmt.Errorf("HTTP %d: %s", resp.StatusCode, body))
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, body)
	}

	// Check for Real-Debrid-style error in JSON
	var rdResp struct {
		Error     string `json:"error"`
		ErrorCode int    `json:"error_code"`
	}
	if err := json.Unmarshal(body, &rdResp); err == nil && rdResp.ErrorCode != 0 {
		return nil, fmt.Errorf("realdebrid error: code=%d, msg='%s'", rdResp.ErrorCode, rdResp.Error)
	}

	return body, nil
}

func PostForm(urlStr string, data url.Values) ([]byte, error) {
	req, err := NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return DoRequest(req)
}
