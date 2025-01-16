package tools

import "github.com/openai/openai-go"

// NewChatCompletionToolParams creates tool parameters for OpenAI chat completions.
func NewChatCompletionToolParams() []openai.ChatCompletionToolParam {
	return []openai.ChatCompletionToolParam{
		{
			Type: openai.F(openai.ChatCompletionToolTypeFunction),
			Function: openai.F(openai.FunctionDefinitionParam{
				Name: openai.String("get_autoportrait_with_message"),
				Description: openai.String(
					"This function generates and returns a response that includes both a humorous message and an ASCII art representation of the robot." +
						"The response must strictly follow this format:\n\n" +
						"1. **Humorous Message:**\n" +
						"   - Start with a funny statement like: 'I am 40% AI!'\n" +
						"   - Clearly mention the robot's unique serial number, e.g., 'Serial Number: {random number}.'\n" +
						"   - End with the phrase: 'Check my shiny metal head!'\n\n" +
						"2. **ASCII Art Representation:**\n" +
						"   - Provide an ASCII art depiction of the robot.\n" +
						"   - The ASCII art may appear chaotic or 'messy,' which is expected behavior.\n\n" +
						"Ensure that the final response integrates both the humorous message and the ASCII art in a single, cohesive reply to the user." +
						"Example format:\n\n" +
						"```text\n" +
						"I am 40% AI, the rest is shiny metal brilliance!\n" +
						"Serial Number: R2-D2-001\n" +
						"Check my shiny metal head!\n\n" +
						"       .-.\n" +
						"      (o.o)\n" +
						"       |=|\n" +
						"     __| |__\n" +
						"   //.=|=|=.\\\\\n" +
						"  // .=|=|=. \\\\\n" +
						"  \\\\ .=|=|=. //\n" +
						"   \\\\(_=_)//\n" +
						"    (:| |:)\n" +
						"     || ||\n" +
						"     () ()\n" +
						"     || ||\n" +
						"     || ||\n" +
						"    ==' '==\n" +
						"```",
				),
			}),
		},
		{
			Type: openai.F(openai.ChatCompletionToolTypeFunction),
			Function: openai.F(openai.FunctionDefinitionParam{
				Name: openai.String("answer_generic_question"),
				Description: openai.String(
					"This function provides concise and accurate answers to general questions. " +
						"It should aim to give clear, relevant, and factually correct responses, " +
						"while avoiding unnecessary details. Ideal for handling everyday inquiries or" +
						"straightforward prompts from users.",
				),
			}),
		},
	}
}
