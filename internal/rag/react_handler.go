package rag

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rag-poc/internal/openaiclient"

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

func (h *OpenAIHandler) reActHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Prompt string `json:"prompt"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	ctx := context.Background()

	question := "What is the weather in New York City?"

	print("> ")
	println(question)
	// Initial parameters for the OpenAI chat completion
	params := openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(question),
		}),
		Tools: openai.F([]openai.ChatCompletionToolParam{
			{
				Type: openai.F(openai.ChatCompletionToolTypeFunction),
				Function: openai.F(openai.FunctionDefinitionParam{
					Name:        openai.String("get_weather"),
					Description: openai.String("Get weather at the given location"),
					Parameters: openai.F(openai.FunctionParameters{
						"type": "object",
						"properties": map[string]interface{}{
							"location": map[string]string{
								"type": "string",
							},
						},
						"required": []string{"location"},
					}),
				}),
			},
		}),
		Seed:  openai.Int(0),
		Model: openai.F(openai.ChatModelGPT4oMini),
	}

	// Make initial chat completion request
	completion, err := h.Client.ChatCompletion(ctx, params)
	if err != nil {
		panic(err)
	}

	toolCalls := completion.Choices[0].Message.ToolCalls

	// Abort early if there are no tool calls
	if len(toolCalls) == 0 {
		fmt.Printf("No function call")
		return
	}

	// If there is a was a function call, continue the conversation
	params.Messages.Value = append(params.Messages.Value, completion.Choices[0].Message)
	for _, toolCall := range toolCalls {
		if toolCall.Function.Name == "get_weather" {
			// Extract the location from the function call arguments
			var args map[string]interface{}
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
				panic(err)
			}
			location := args["location"].(string)

			// Simulate getting weather data
			weatherData := getWeather(location)

			// Print the weather data
			fmt.Printf("Weather in %s: %s\n", location, weatherData)

			params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(toolCall.ID, weatherData))
		}
	}

	completion, err = h.Client.ChatCompletion(ctx, params)
	if err != nil {
		panic(err)
	}

	println(completion.Choices[0].Message.Content)
}

// Mock function to simulate weather data retrieval
func getWeather(location string) string {
	// In a real implementation, this function would call a weather API
	return fmt.Sprintf("In %s, it is Sunny, 25°C", location)
}
