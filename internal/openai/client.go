package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"

	"github.com/techcorrectco/reqd/internal"
)

type ValidationResponse struct {
	Input       string   `json:"input"`
	Problems    []string `json:"problems"`
	Recommended string   `json:"recommended"`
	Keyword     string   `json:"keyword"`
}

type OpenAIRequest struct {
	Model          string            `json:"model"`
	Messages       []Message         `json:"messages"`
	ResponseFormat map[string]string `json:"response_format,omitempty"`
	Temperature    float64           `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

func ValidateRequirement(input string) (*ValidationResponse, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set - use --no-validate to skip validation")
	}

	// Parse the prompt template
	tmpl, err := template.New("prompt").Parse(internal.ValidateRequirementPrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse prompt template: %w", err)
	}

	// Execute template with input
	var promptBuffer bytes.Buffer
	err = tmpl.Execute(&promptBuffer, map[string]string{"Input": input})
	if err != nil {
		return nil, fmt.Errorf("failed to execute prompt template: %w", err)
	}

	// Create OpenAI request
	reqBody := OpenAIRequest{
		Model: "gpt-4o",
		Messages: []Message{
			{
				Role:    "user",
				Content: promptBuffer.String(),
			},
		},
		ResponseFormat: map[string]string{
			"type": "json_object",
		},
		Temperature: 0.0,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make HTTP request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse OpenAI response
	var openaiResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openaiResp); err != nil {
		return nil, fmt.Errorf("failed to decode OpenAI response: %w", err)
	}

	if len(openaiResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in OpenAI response")
	}

	// Parse the JSON content from OpenAI
	responseContent := openaiResp.Choices[0].Message.Content
	var validationResp ValidationResponse
	if err := json.Unmarshal([]byte(responseContent), &validationResp); err != nil {
		return nil, fmt.Errorf("failed to parse validation response JSON: %w", err)
	}

	return &validationResp, nil
}