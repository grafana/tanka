package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/grafana/tanka/pkg/agent/tools"
)

const geminiAPIURLTemplate = "https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s"

// GeminiProvider calls the Google Gemini API directly over HTTP.
type GeminiProvider struct {
	apiKey string
	model  string
	client *http.Client
}

// NewGeminiProvider creates a provider that uses the Gemini API.
func NewGeminiProvider(apiKey, model string) *GeminiProvider {
	return &GeminiProvider{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{},
	}
}

func (p *GeminiProvider) Name() string { return ProviderGemini }

// geminiPart is a part of a Gemini content block.
type geminiPart struct {
	Text             string                  `json:"text,omitempty"`
	FunctionCall     *geminiFunctionCall     `json:"functionCall,omitempty"`
	FunctionResponse *geminiFunctionResponse `json:"functionResponse,omitempty"`
}

type geminiFunctionCall struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
}

type geminiFunctionResponse struct {
	Name     string                 `json:"name"`
	Response map[string]interface{} `json:"response"`
}

// geminiContent is a message in the Gemini conversation format.
type geminiContent struct {
	Role  string       `json:"role"`
	Parts []geminiPart `json:"parts"`
}

// geminiTool is a tool definition for the Gemini API.
type geminiTool struct {
	FunctionDeclarations []geminiFunctionDecl `json:"functionDeclarations"`
}

type geminiFunctionDecl struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"`
}

// geminiRequest is the request body for the Gemini generateContent API.
type geminiRequest struct {
	SystemInstruction *geminiContent  `json:"systemInstruction,omitempty"`
	Contents          []geminiContent `json:"contents"`
	Tools             []geminiTool    `json:"tools,omitempty"`
}

// geminiResponse is the response from the Gemini generateContent API.
type geminiResponse struct {
	Candidates []struct {
		Content      geminiContent `json:"content"`
		FinishReason string        `json:"finishReason"`
	} `json:"candidates"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error,omitempty"`
}

func (p *GeminiProvider) Chat(ctx context.Context, systemPrompt string, messages []Message, toolDefs []tools.Tool) (*Message, error) {
	// Convert messages to Gemini format
	geminiContents := make([]geminiContent, 0, len(messages))
	for _, msg := range messages {
		var parts []geminiPart
		role := "user"
		if msg.Role == RoleAssistant {
			role = "model"
		}

		for _, c := range msg.Content {
			switch c.Type {
			case ContentTypeText:
				parts = append(parts, geminiPart{Text: c.Text})
			case ContentTypeToolUse:
				var args map[string]interface{}
				_ = json.Unmarshal(c.Input, &args)
				parts = append(parts, geminiPart{
					FunctionCall: &geminiFunctionCall{
						Name: c.Name,
						Args: args,
					},
				})
			case ContentTypeToolResult:
				// Tool results in Gemini go in a "user" role with functionResponse
				result := map[string]interface{}{"result": c.Text}
				if c.IsError {
					result["error"] = c.Text
					delete(result, "result")
				}
				parts = append(parts, geminiPart{
					FunctionResponse: &geminiFunctionResponse{
						Name:     c.Name,
						Response: result,
					},
				})
			}
		}
		geminiContents = append(geminiContents, geminiContent{
			Role:  role,
			Parts: parts,
		})
	}

	// Build tool declarations
	var geminiTools []geminiTool
	if len(toolDefs) > 0 {
		decls := make([]geminiFunctionDecl, 0, len(toolDefs))
		for _, t := range toolDefs {
			decls = append(decls, geminiFunctionDecl{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.Schema,
			})
		}
		geminiTools = []geminiTool{{FunctionDeclarations: decls}}
	}

	reqBody := geminiRequest{
		SystemInstruction: &geminiContent{
			Parts: []geminiPart{{Text: systemPrompt}},
		},
		Contents: geminiContents,
		Tools:    geminiTools,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshalling request: %w", err)
	}

	apiURL := fmt.Sprintf(geminiAPIURLTemplate, p.model, p.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("content-type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling Gemini API: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp geminiResponse
		if json.Unmarshal(respBytes, &errResp) == nil && errResp.Error != nil {
			return nil, fmt.Errorf("gemini API error %d %s: %s", errResp.Error.Code, errResp.Error.Status, errResp.Error.Message)
		}
		return nil, fmt.Errorf("gemini API returned HTTP %d: %s", resp.StatusCode, string(respBytes))
	}

	var apiResp geminiResponse
	if err := json.Unmarshal(respBytes, &apiResp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	if len(apiResp.Candidates) == 0 {
		return nil, fmt.Errorf("gemini returned no candidates")
	}

	// Convert response to internal format
	result := &Message{Role: RoleAssistant}
	for _, part := range apiResp.Candidates[0].Content.Parts {
		switch {
		case part.FunctionCall != nil:
			inputBytes, _ := json.Marshal(part.FunctionCall.Args)
			// Generate a synthetic ID for the tool call
			id := fmt.Sprintf("call_%s", strings.ReplaceAll(part.FunctionCall.Name, "_", ""))
			result.Content = append(result.Content, Content{
				Type:  ContentTypeToolUse,
				ID:    id,
				Name:  part.FunctionCall.Name,
				Input: inputBytes,
			})
		case part.Text != "":
			result.Content = append(result.Content, Content{
				Type: ContentTypeText,
				Text: part.Text,
			})
		}
	}

	return result, nil
}
