package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rag-poc/internal/openaiclient"
	"rag-poc/internal/rag/tools"
	"rag-poc/internal/rag/utils"

	"github.com/openai/openai-go"
)

type OpenAIHandler struct {
	Client openaiclient.Service
}

func NewHandler(client openaiclient.Service) *OpenAIHandler {
	return &OpenAIHandler{
		Client: client,
	}
}

func (h *OpenAIHandler) ReActHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Prompt string `json:"prompt"`
	}

	if err := utils.ParseJSON(r, &req); err != nil {
		log.Printf("Error decoding request: %v", err)
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}
	ctx := context.Background()
	response, err := h.handleChatCompletion(ctx, req.Prompt)
	if err != nil {
		log.Printf("Error handling chat completion: %v", err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}

	utils.SendJSONResponse(w, http.StatusCreated, map[string]string{"response": response})
}

// handleChatCompletion orchestrates the OpenAI chat completion logic
func (h *OpenAIHandler) handleChatCompletion(ctx context.Context, prompt string) (string, error) {
	toolParams := tools.NewChatCompletionToolParams() // Reuse tools definition
	params := openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		}),
		Tools:       openai.F(toolParams),
		Model:       openai.F(openai.ChatModelGPT4oMini),
		Temperature: openai.Float(0),
	}

	completion, err := h.Client.ChatCompletion(ctx, params)
	if err != nil {
		return "", err
	}

	toolCalls := completion.Choices[0].Message.ToolCalls
	if len(toolCalls) == 0 {
		return "No function call", nil
	}

	params.Messages.Value = append(params.Messages.Value, completion.Choices[0].Message)

	// Pass params by reference to modify directly
	err = h.handleToolCalls(toolCalls, &params) // Pass pointer here
	if err != nil {
		return "", err
	}

	completion, err = h.Client.ChatCompletion(ctx, params)
	if err != nil {
		return "", err
	}

	return completion.Choices[0].Message.Content, nil
}

// handleToolCalls processes tool calls and updates messages
func (h *OpenAIHandler) handleToolCalls(toolCalls []openai.ChatCompletionMessageToolCall, params *openai.ChatCompletionNewParams) error {
	for _, toolCall := range toolCalls {
		switch toolCall.Function.Name {
		case tools.NewChatCompletionToolParams()[0].Function.Value.Name.Value:
			var args map[string]interface{}
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
				return err
			}

			location := args["location"].(string)
			weatherData := tools.GetWeather(location)

			// Append the tool message response for this tool call ID
			params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(toolCall.ID, weatherData))

		default:
			// Handle unknown or unsupported function calls
			log.Printf("Unhandled tool function: %s", toolCall.Function.Name)
			params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(toolCall.ID, fmt.Sprintf("Function '%s' is not supported.", toolCall.Function.Name)))
		}
	}
	return nil
}
