package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
	"vovan_ai/config"
)

type AiClient struct {
	BaseUlr   string
	ApiKey    string
	ModelName string
	Client    *http.Client
}

func NewAiClient(cnf *config.Config) *AiClient {
	return &AiClient{
		BaseUlr:   cnf.BaseUrl,
		ApiKey:    cnf.ApiKey,
		ModelName: cnf.ModelName,
		Client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

type RequestPayload struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

type StreamResponse struct {
	Choices []StreamChoice `json:"choices"`
}

type StreamChoice struct {
	Delta Delta `json:"delta"`
}

type Delta struct {
	Content string `json:"content"`
}

func (c *AiClient) SendChat(messages []Message) (string, error) {
	req, err := c.createRequest(messages, false)
	if err != nil {
		return "", err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result Response

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", errors.New("не удалось получить ответ")
	}

	return result.Choices[0].Message.Content, nil
}

func (c *AiClient) SendStreamChat(messages []Message, onChunk func(string)) (string, error) {
	req, err := c.createRequest(messages, true)
	if err != nil {
		return "", err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result strings.Builder

	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")

		if data == "[DONE]" {
			break
		}

		var chunk StreamResponse
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}

		if len(chunk.Choices) == 0 {
			continue
		}

		content := chunk.Choices[0].Delta.Content

		result.WriteString(content)

		onChunk(content)
	}

	return result.String(), nil
}

func (c *AiClient) createRequest(messages []Message, stream bool) (*http.Request, error) {
	payload := RequestPayload{
		Model:    c.ModelName,
		Messages: messages,
		Stream:   stream,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		c.BaseUlr+"/chat/completions",
		bytes.NewReader(data),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.ApiKey)
	req.Header.Set("Content-type", "application/json")

	return req, nil
}
