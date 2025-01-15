package handler

import (
	"context"
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
	toolParams := tools.NewChatCompletionToolParams()
	params := openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are an assistant capable of generating and presenting your self-portrait as ASCII art."),
			openai.UserMessage(prompt),
		}),
		Tools:       openai.F(toolParams),
		Model:       openai.F(openai.ChatModelGPT3_5Turbo),
		Temperature: openai.Float(0),
	}

	// Make initial chat completion request
	completion, err := h.Client.ChatCompletion(ctx, params)
	if err != nil {
		panic(err)
	}

	toolCalls := completion.Choices[0].Message.ToolCalls
	params.Messages.Value = append(params.Messages.Value, completion.Choices[0].Message)

	for _, toolCall := range toolCalls {
		if toolCall.Function.Name == "get_autoportrait" {
			asciiArt := tools.GetAutoportrait()

			params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(toolCall.ID, asciiArt))
		} else {
			log.Printf("Unhandled tool call: %s", toolCall.Function.Name)
		}
	}
	// Secondary completion
	completion, err = h.Client.ChatCompletion(ctx, params)
	if err != nil {
		return "", fmt.Errorf("secondary chat completion failed: %w", err)
	}

	if len(completion.Choices) == 0 || completion.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("no valid response received in secondary completion")
	}
	return completion.Choices[0].Message.Content, nil
}
