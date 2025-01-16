package handler

import (
	"context"
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
			openai.SystemMessage("You are an assistant capable of help with various of tasks."),
			openai.UserMessage(prompt),
		}),
		Tools:       openai.F(toolParams),
		Model:       openai.F(openai.ChatModelGPT4oMini),
		Temperature: openai.Float(0),
	}

	// Make initial chat completion request
	initCompletion, err := h.Client.ChatCompletion(ctx, params)
	if err != nil {
		panic(err)
	}

	toolCalls := initCompletion.Choices[0].Message.ToolCalls
	params.Messages.Value = append(params.Messages.Value, initCompletion.Choices[0].Message)

	for _, toolCall := range toolCalls {
		if toolCall.Function.Name == "get_autoportrait_with_message" {
			asciiArt, err := tools.GetAutoportrait()
			if err != nil {
				log.Printf("Error fetching autoportrait: %v", err)
				asciiArt = "Sorry, I couldn't fetch the autoportrait at the moment. Please try again later!"
			}

			params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(toolCall.ID, asciiArt))
		} else if toolCall.Function.Name == "answer_generic_questions" {
			answer, err := h.Client.ChatCompletion(ctx, params)
			if err != nil {
				log.Printf("Error processing generic question: %v", err)
				fallbackMessage := "Apologies, I couldn't process your question right now. Please try again later!"
				params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(toolCall.ID, fallbackMessage))
			} else {
				params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(toolCall.ID, answer.Choices[0].Message.Content))
			}
		} else {
			log.Printf("Unhandled tool call: %s", toolCall.Function.Name)
			params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(toolCall.ID, "Sorry, I don't know how to handle this request yet."))
		}
	}
	// Secondary completion
	secondCompletion, err := h.Client.ChatCompletion(ctx, params)
	if err != nil {
		return secondCompletion.Choices[0].Message.Content, nil
	}

	// if len(completion.Choices) == 0 || completion.Choices[0].Message.Content == "" {
	// 	return "", fmt.Errorf("no valid response received in secondary completion")
	// }
	return secondCompletion.Choices[0].Message.Content, nil
}
