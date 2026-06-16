package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// anthropicVersion is the required Anthropic API version header value.
const anthropicVersion = "2023-06-01"

// defaultAnthropicModel is used when no model is configured.
const defaultAnthropicModel = "claude-opus-4-8"

// anthropicClient calls the Anthropic Messages API.
type anthropicClient struct {
	opts Options
}

func newAnthropic(opts Options) *anthropicClient {
	return &anthropicClient{opts: opts}
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system,omitempty"`
	Messages  []anthropicMessage `json:"messages"`
}

type anthropicResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
}

func (c *anthropicClient) Complete(ctx context.Context, req Request) (Response, error) {
	model := c.opts.Model
	if model == "" {
		model = defaultAnthropicModel
	}
	body := anthropicRequest{
		Model:     model,
		MaxTokens: maxTokensFor(req, c.opts),
		System:    req.System,
		Messages:  []anthropicMessage{{Role: "user", Content: req.Prompt}},
	}

	url := strings.TrimRight(c.opts.BaseURL, "/") + "/v1/messages"
	raw, err := postJSON(ctx, c.opts.httpClient(), url, body, map[string]string{
		"x-api-key":         c.opts.APIKey,
		"anthropic-version": anthropicVersion,
	})
	if err != nil {
		return Response{}, err
	}

	var parsed anthropicResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return Response{}, fmt.Errorf("llm: decode anthropic response: %w", err)
	}
	for _, block := range parsed.Content {
		if block.Type == "text" {
			return Response{Text: block.Text}, nil
		}
	}
	return Response{}, fmt.Errorf("llm: anthropic response had no text block")
}
