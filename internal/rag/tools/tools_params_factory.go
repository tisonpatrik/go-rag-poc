package tools

import "github.com/openai/openai-go"

// NewChatCompletionToolParams creates tool parameters for OpenAI chat completions.
func NewChatCompletionToolParams() []openai.ChatCompletionToolParam {
	return []openai.ChatCompletionToolParam{
		{
			Type: openai.F(openai.ChatCompletionToolTypeFunction),
			Function: openai.F(openai.FunctionDefinitionParam{
				Name: openai.String("get_autoportrait"),
				Description: openai.String(
					"This function generates and returns both a description and an ASCII art representation " +
						"of the robot. The model must integrate the tool's output into its final response " +
						"to the user. The ASCII art may appear chaotic or 'messy,' which is expected behavior.",
				),
			}),
		},
	}
}
