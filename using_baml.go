package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"

	b "agent/baml_client"
	"agent/baml_client/types"
)

func UsingBAML() {
	scanner := bufio.NewScanner(os.Stdin)
		
	// Initial prompt is defined here:  baml_src/chat_with_tool.baml 33
	var conversation []types.Message

	userLoop: for {
		print_(BLUE, "You")
		input, ok := getUserMessage(scanner); if !ok { break; }

		msg := types.Message{ Role: "user", Msg: input }
		conversation = append(conversation, msg)

		for {
			// Call AI using BAML -- messages are combined into a single prompt with the __full__ conversation history
			aiResponse, err := b.CallGPT5Mini(context.Background(), conversation)
			if err != nil { fmt.Printf("Error: %s\n", err.Error()); break userLoop; }

			aiResponseBytes, _ := json.Marshal(aiResponse)
			aiResponseString := string(aiResponseBytes)

			conversation = append(conversation, types.Message{ Role: "assistant", Msg: aiResponseString })

			// Response is already structured
			toolResult, wasUsed := useTool_STRUCTURED(aiResponse)

			if wasUsed {
				print(GREEN, "Tool:in", aiResponseString)
				print(GREEN, "Tool:out", toolResult)
				conversation = append(conversation, types.Message{ Role: "tool-response", Msg: toolResult })
			} else {
				print(YELLOW, "AI", aiResponse.AsTextualReply().Value)
				break // Continue the user loop
			}
		}
	}
}
