package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"pingpong-backend/internal/config"
)

var ErrAgentDisabled = errors.New("agent is disabled")
var ErrLLMNotConfigured = errors.New("llm api key is not configured")

type LLMClient struct {
	baseURL    string
	apiKey     string
	model      string
	httpClient *http.Client
}

type chatCompletionRequest struct {
	Model       string           `json:"model"`
	Messages    []chatMessage    `json:"messages"`
	Temperature float64          `json:"temperature"`
	Response    *responseOptions `json:"response_format,omitempty"`
}

type responseOptions struct {
	Type string `json:"type"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

func NewLLMClient(cfg config.Config) (*LLMClient, error) {
	if !cfg.AgentEnabled {
		return nil, ErrAgentDisabled
	}
	if cfg.LLMAPIKey == "" {
		return nil, ErrLLMNotConfigured
	}

	timeout := time.Duration(cfg.AgentTimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	return &LLMClient{
		baseURL: strings.TrimRight(cfg.LLMBaseURL, "/"),
		apiKey:  cfg.LLMAPIKey,
		model:   cfg.LLMModel,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}, nil
}

func (c *LLMClient) Chat(ctx context.Context, systemPrompt, userPrompt string, jsonResponse bool) (string, error) {
	request := chatCompletionRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.2,
	}
	if jsonResponse {
		request.Response = &responseOptions{Type: "json_object"}
	}

	body, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	httpRequest.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpRequest.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		var completion chatCompletionResponse
		if err := json.Unmarshal(responseBody, &completion); err == nil && completion.Error != nil && completion.Error.Message != "" {
			return "", fmt.Errorf("llm request failed: %s", completion.Error.Message)
		}
		message := strings.TrimSpace(string(responseBody))
		if message != "" {
			return "", fmt.Errorf("llm request failed with status %d: %s", response.StatusCode, truncateMessage(message, 240))
		}
		return "", fmt.Errorf("llm request failed with status %d", response.StatusCode)
	}

	var completion chatCompletionResponse
	if err := json.Unmarshal(responseBody, &completion); err != nil {
		return "", fmt.Errorf("decode llm response: %w", err)
	}
	if completion.Error != nil {
		if completion.Error != nil && completion.Error.Message != "" {
			return "", fmt.Errorf("llm request failed: %s", completion.Error.Message)
		}
		return "", errors.New("llm request failed")
	}
	if len(completion.Choices) == 0 {
		return "", errors.New("llm returned no choices")
	}
	return completion.Choices[0].Message.Content, nil
}

func truncateMessage(value string, maxLength int) string {
	if len(value) <= maxLength {
		return value
	}
	return value[:maxLength] + "..."
}
