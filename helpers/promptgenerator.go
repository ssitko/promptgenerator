package promptgenerator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type PromptMode int

type Modes struct {
	Explanations  bool
	Comments      bool
	Documentation bool
}

const (
	BasicPrompt PromptMode = iota
	ExplanationsPrompt
	CommentsPrompt
	DocumentationPrompt
)

type RequestBody struct {
	Contents         []Content        `json:"contents"`
	SafetySettings   []SafetySetting  `json:"safetySettings,omitempty"`
	GenerationConfig GenerationConfig `json:"generationConfig"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type SafetySetting struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

type GenerationConfig struct {
	StopSequences   []string `json:"stopSequences"`
	Temperature     float64  `json:"temperature"`
	MaxOutputTokens int      `json:"maxOutputTokens"`
	TopP            float64  `json:"topP"`
	TopK            int      `json:"topK"`
}

type ResponseBody struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content struct {
		Parts []Part `json:"parts"`
	} `json:"content"`
}

// Implementation of PromptHandler
type PromptHandler struct {
	url             string
	temperature     float64
	topP            float64
	topK            int
	maxOutputTokens int
}

func NewPromptHandler(baseUrl, apiKey string, temperature, topP float64, maxOutputTokens, topK int) *PromptHandler {
	return &PromptHandler{
		url:             fmt.Sprintf("%s%s", baseUrl, apiKey),
		temperature:     temperature,
		maxOutputTokens: maxOutputTokens,
		topP:            topP,
		topK:            topK,
	}
}

func (p *PromptHandler) GenerateContent(actor, prompt string, modes Modes) (string, error) {
	instructions := ""
	if modes.Comments {
		instructions = fmt.Sprintf("Include detailed in-code comments whenever applicable.\n%s", instructions)
	}
	if modes.Documentation {
		instructions = fmt.Sprintf("Attach extensive documentation for generated code as well as for used dependencies or modules. Include it below code and DO comment these lines out.\n%s", instructions)
	}
	if modes.Explanations {
		instructions = fmt.Sprintf("Add detailed explanation on how to use code and what it does. Include them below code and DO comment these lines out.\n%s", instructions)
	}
	prompt = fmt.Sprintf("%s.\n\n%s", prompt, instructions)

	requestBody := RequestBody{
		Contents: []Content{{
			Parts: []Part{{Text: fmt.Sprintf("As a %s, %s", actor, prompt)}},
		}},
		GenerationConfig: GenerationConfig{
			StopSequences:   []string{},
			Temperature:     p.temperature,
			MaxOutputTokens: p.maxOutputTokens,
			TopP:            p.topP,
			TopK:            p.topK,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, p.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var responseBody ResponseBody
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return "", err
	}

	// Ensure response contains valid data
	if len(responseBody.Candidates) == 0 || len(responseBody.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no valid text found in response")
	}

	return removeFirstLine(responseBody.Candidates[0].Content.Parts[0].Text), nil
}

func removeFirstLine(code string) string {
	lines := strings.Split(code, "\n")
	if len(lines) > 1 {
		return strings.Join(lines[1:], "\n")
	}
	return ""
}
