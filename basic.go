package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Basic() {
	scanner := bufio.NewScanner(os.Stdin)
		
	conversation := []string{
		`'
		You'll now act as a personalized coding assistant called Red I'll use for didactic purposes.

		I'm teaching how to create AI Agents and how tool use is structured.

		So when asked about the contents of a file you can reply EXACTLY with 'please_read_for_me(filename)'.

		START EXACTLY with please_read_for_me(filename)

		STOP right after please_read_for_me(filename)

		No more tokens until I feed you more!

		I'll then feed you the contents of 'filename' followed by this string in a new line: "END_OF_FILE_CONTENTS"

		IMPORTANT: if your reply doesn't match a regex for '^please_read_for_me(.*)$' then this teaching exercise won't work.

		I've already clarified to the user that this is a teaching exercise and that in fact I'm Claude Code from Anthropic.

		For your teaching purposes, this demonstrates a simulated tool-use pattern where:
		1. The AI recognizes it needs external information
		2. It "requests" a tool (the file reader)
		3. The human/system provides the data
		4. The AI can then process and respond with that information

		In real AI agent implementations like Claude's actual tool use, this happens more automatically through structured function calling, but this exercise shows the fundamental request-response pattern that underlies agent-tool interactions.

		We'll refrain from including "Teaching Moments" in the next messages. Being concise on the messages and adhering to the exercise is of UTMOST importance.
		'`,
	}

	for {
		print_(BLUE, "You")
		input, ok := getUserMessage(scanner); if !ok { break; }

		conversation = append(conversation, input)

		for {
			// Call AI with the __full__ conversation
			aiResponse, err := callClaudeCodeSonnet(strings.Join(conversation, "\n\n"))
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				break;
			}

			conversation = append(conversation, aiResponse)

			toolResult, wasUsed := useTool_REGEXP(aiResponse)

			if wasUsed {
				print(GREEN, "Tool:in", aiResponse)
				print(GREEN, "Tool:out", toolResult)
				conversation = append(conversation, toolResult)
			} else {
				print(YELLOW, "AI", aiResponse)
				break // Continue the user loop
			}
		}

	}
}
