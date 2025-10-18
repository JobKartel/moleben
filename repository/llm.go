package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

type Client struct {
	HTTP     *http.Client
	BaseURL  string
	APIKey   string
	Model    string
	Provider string
	Referer  string
	Title    string
}

func NewClient(baseURL, apiKey, model, provider, referer, title string) *Client {
	return &Client{
		HTTP:     &http.Client{Timeout: 60 * time.Second},
		BaseURL:  baseURL,
		APIKey:   apiKey,
		Model:    model,
		Provider: provider,
		Referer:  referer,
		Title:    title,
	}
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

var ErrEmptyResponse = errors.New("llm returned empty content")

func (c *Client) Chat(ctx context.Context, systemPrompt, user string) (string, error) {
	req := chatRequest{Model: c.Model, Messages: []chatMessage{{"system", systemPrompt}, {"user", user}}, Stream: false}
	b, _ := json.Marshal(req)
	endpoint := c.BaseURL + "/chat/completions"
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(b))
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	if c.Provider == "openrouter" {
		if c.Referer != "" {
			httpReq.Header.Set("HTTP-Referer", c.Referer)
		}
		if c.Title != "" {
			httpReq.Header.Set("X-Title", c.Title)
		}
	}
	resp, err := c.HTTP.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", errors.New("llm http error: " + resp.Status + ": " + string(body))
	}
	var ch chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&ch); err != nil {
		return "", err
	}
	if len(ch.Choices) == 0 || ch.Choices[0].Message.Content == "" {
		return "", ErrEmptyResponse
	}
	return ch.Choices[0].Message.Content, nil
}
