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
	"github.com/techcorrectco/reqd/internal/types"
)

type ValidationResponse struct {
	Input       string   `json:"input"`
	Problems    []string `json:"problems"`
	Recommended string   `json:"recommended"`
}

type ParentProposalResponse struct {
	ProposedParent *string `json:"proposed_parent"`
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

// renderTemplate renders a template string with provided data
func renderTemplate(templateStr string, data map[string]string) (string, error) {
	tmpl, err := template.New("prompt").Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse prompt template: %w", err)
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute prompt template: %w", err)
	}

	return buffer.String(), nil
}

// makeOpenAIRequest sends a request to OpenAI and returns the response content
func makeOpenAIRequest(prompt string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	// Create OpenAI request
	reqBody := OpenAIRequest{
		Model: "gpt-4o",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		ResponseFormat: map[string]string{
			"type": "json_object",
		},
		Temperature: 0.0,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make HTTP request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse OpenAI response
	var openaiResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openaiResp); err != nil {
		return "", fmt.Errorf("failed to decode OpenAI response: %w", err)
	}

	if len(openaiResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in OpenAI response")
	}

	return openaiResp.Choices[0].Message.Content, nil
}

func ValidateRequirement(input string) (*ValidationResponse, error) {
	// Render template
	prompt, err := renderTemplate(internal.ValidateRequirementPrompt, map[string]string{"Input": input})
	if err != nil {
		return nil, err
	}

	// Make OpenAI request
	responseContent, err := makeOpenAIRequest(prompt)
	if err != nil {
		return nil, err
	}

	// Parse the JSON content from OpenAI
	var validationResp ValidationResponse
	if err := json.Unmarshal([]byte(responseContent), &validationResp); err != nil {
		return nil, fmt.Errorf("failed to parse validation response JSON: %w", err)
	}

	return &validationResp, nil
}

func ProposeParent(requirement string, branches []types.Requirement) (*ParentProposalResponse, error) {
	// Format branches using DisplayFormat method
	var parentsText string
	for _, branch := range branches {
		parentsText += branch.DisplayFormat() + "\n"
	}

	// Render template
	prompt, err := renderTemplate(internal.ProposeParentPrompt, map[string]string{
		"Parents":     parentsText,
		"Requirement": requirement,
	})
	if err != nil {
		return nil, err
	}

	// Make OpenAI request
	responseContent, err := makeOpenAIRequest(prompt)
	if err != nil {
		return nil, err
	}

	// Parse the JSON content from OpenAI
	var proposalResp ParentProposalResponse
	if err := json.Unmarshal([]byte(responseContent), &proposalResp); err != nil {
		return nil, fmt.Errorf("failed to parse parent proposal response JSON: %w", err)
	}

	return &proposalResp, nil
}
