package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"agent/baml_client/types"
)

func main() {
	fmt.Print("\n\n\n")

	// print(RED, "LOOP", "Basic()")
	// Basic()
	// print(RED, "LOOP", "UsingBAML()")
	// UsingBAML()
	print(RED, "LOOP", "UsingBAMLParser()")
	UsingBAMLParser()
}

const COLOR_CLEAR = "\u001b[0m"
const BLUE = "\u001b[94m"
const YELLOW = "\u001b[93m"
const GREEN = "\u001b[92m"
const RED = "\u001b[91m"

func print(color, name string, printfArgs ...any) {
	msg := ""
	if len(printfArgs) > 0 {
		msg = fmt.Sprintf(printfArgs[0].(string), printfArgs[1:]...)
	}
	fmt.Printf("%s%s%s: %s\n", color, name, COLOR_CLEAR, msg)
}

func print_(color, name string, printfArgs ...any) {
	msg := ""
	if len(printfArgs) > 0 {
		msg = fmt.Sprintf(printfArgs[0].(string), printfArgs[1:]...)
	}
	fmt.Printf("%s%s%s: %s", color, name, COLOR_CLEAR, msg)
}

func getUserMessage (scanner *bufio.Scanner) (string, bool) {
	if !scanner.Scan() {
		return "", false
	}
	return scanner.Text(), true
}

func useTool_REGEXP(aiResponse string) (string, bool) {
	readFileRegex := regexp.MustCompile(`^please_read_for_me\((.+?)\)$`)

	cleanedResponse := strings.TrimSpace(aiResponse)

	readFileMatch := readFileRegex.FindStringSubmatch(cleanedResponse)

	if readFileMatch != nil {
		filename := readFileMatch[1]
		data, err := os.ReadFile(filename)

		if err != nil {
			// Let the AI know about the error
			return err.Error(), true
		}

		readFileOutput := fmt.Sprintf("Here's the content of %s:\n%s\nEND_OF_FILE_CONTENTS", filename, string(data))

		return readFileOutput, true
	}

	return "", false
}

func useTool_STRUCTURED(aiResponse types.Union2TextualReplyOrToolUse) (string, bool) {
	if !aiResponse.IsToolUse() { return "", false }

	toolUse := aiResponse.AsToolUse()

	if toolUse.ToolName == "please_read_for_me" {
		data, err := os.ReadFile(toolUse.Filename)

		if err != nil {
			// Let the AI know about the error
			return err.Error(), true
		}

		readFileOutput := fmt.Sprintf(`{ filename: %s, contents: %s }`, toolUse.Filename, string(data))

		return readFileOutput, true
	}

	return "", false
}

func callClaudeCodeSonnet(input string) (string, error) {
	// cmd := exec.Command("claude", "--tools", "''", "-p", "--model", "sonnet", "--system-prompt", SYSTEM_PROMPT, "'" + input + "'")
	cmd := exec.Command("claude", "--tools", "''", "-p", "--model", "sonnet", "'" + input + "'")

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	cmd.Stdin = nil

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("getAiResponse failed: %w\nstderr: %s\n stdout: %s", err, stderr.String(), stdout.String())
	}

	return stdout.String(), nil
}
