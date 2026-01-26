package main

import (
	"bufio"
	"fmt"
	"os"

	b "agent/baml_client"
	"agent/baml_client/types"
)

func UsingBAMLParser() {
	scanner := bufio.NewScanner(os.Stdin)
		
	messages := []types.Message{
		{ Role: "system", Msg: `'

			You'll now act as a personalized coding assistant called Red I'll use for didactic purposes.

			I'm teaching how to create AI Agents and how tool use is structured.

			So when asked about the contents of a file you can

			Answer in JSON using this schema to access a tool to read a file

			{
			  type: "ToolUse",
			  toolName: "please_read_for_me",
			  filename: string,
			}

			or use reply with a normal message following this JSON schema: 

			{
			  type: "Text",
			  value: string,
			}

		'`},
	}

	userLoop: for {
		print_(BLUE, "You")
		input, ok := getUserMessage(scanner); if !ok { break; }

		msg := types.Message{ Role: "user", Msg: input }
		messages = append(messages, msg)

		for {
			// Call AI with the __full__ conversation history
			aiResponse, err := callClaudeCodeSonnet(joinConversation(messages))
			if err != nil { fmt.Printf("Error: %s\n", err.Error()); break userLoop; }

			// Then parse with BAML
			parsed, err := b.Parse.CallGPT5Mini(aiResponse)
			if err != nil { fmt.Printf("Msg: %s\nError: %s\n", aiResponse, err.Error()); break userLoop; }

			messages = append(messages, types.Message{ Role: "assistant", Msg: aiResponse })

			toolResult, wasUsed := useTool_STRUCTURED(parsed)

			if wasUsed {
				print(GREEN, "Tool:in", aiResponse)
				print(GREEN, "Tool:out", toolResult)
				messages = append(messages, types.Message{ Role: "tool-response", Msg: toolResult })
			} else {
				print(YELLOW, "AI", parsed.AsTextualReply().Value)
				break // Continue the user loop
			}
		}

	}
}

func joinConversation(theConversation []types.Message) string {
	joined := ""
	for _, message := range theConversation {
		// map and join
		joined += string(message.Role)
		joined += message.Msg
	}
	return joined
}
