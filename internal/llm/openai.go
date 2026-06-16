package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// openAIClient calls an OpenAI-compatible chat completions endpoint.
type openAIClient struct {
	opts Options
}

func newOpenAI(opts Options) *openAIClient {
	return &openAIClient{opts: opts}
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	Messages  []openAIMessage `json:"messages"`
}

type openAIResponse struct {
	Choices []struct {
		Message openAIMessage `json:"message"`
	} `json:"choices"`
}

func (c *openAIClient) Complete(ctx context.Context, req Request) (Response, error) {
	body := openAIRequest{
		Model:     c.opts.Model,
		MaxTokens: maxTokensFor(req, c.opts),
		Messages:  []openAIMessage{},
	}
	if req.System != "" {
		body.Messages = append(body.Messages, openAIMessage{Role: "system", Content: req.System})
	}
	body.Messages = append(body.Messages, openAIMessage{Role: "user", Content: req.Prompt})

	url := strings.TrimRight(c.opts.BaseURL, "/") + "/chat/completions"
	raw, err := postJSON(ctx, c.opts.httpClient(), url, body, map[string]string{
		"Authorization": "Bearer " + c.opts.APIKey,
	})
	if err != nil {
		return Response{}, err
	}

	var parsed openAIResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return Response{}, fmt.Errorf("llm: decode openai response: %w", err)
	}
	if len(parsed.Choices) == 0 {
		return Response{}, fmt.Errorf("llm: openai response had no choices")
	}
	return Response{Text: parsed.Choices[0].Message.Content}, nil
}

// maxTokensFor prefers the per-request cap, then the configured one.
func maxTokensFor(req Request, opts Options) int {
	if req.MaxTokens > 0 {
		return req.MaxTokens
	}
	return opts.maxTokens()
}

// postJSON posts v as JSON to url with the given extra headers and returns the
// response body, or an error for any non-2xx status.
func postJSON(ctx context.Context, client *http.Client, url string, v any, headers map[string]string) ([]byte, error) {
	payload, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	for k, val := range headers {
		httpReq.Header.Set(k, val)
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("llm: provider returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}
	return raw, nil
}
