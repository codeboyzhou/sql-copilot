package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/codeboyzhou/sql-copilot/agentic/llm"
	"github.com/codeboyzhou/sql-copilot/strconst"
)

const MaxLoops = 10

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("%s No .env file found: %v\n", strconst.EmojiError, err)
	}

	model := llm.DeepSeekChat
	fmt.Printf("%s Initializing agent with model: %s\n", strconst.EmojiRunning, model.Name)
	var agent llm.OpenAICompatibleClient
	agent.InitializeWithModel(model)
	fmt.Printf("%s Successfully initialized agent with model: %s\n", strconst.EmojiSuccess, model.Name)

	userPromptHistory := []string{}

	for i := 0; i < MaxLoops; i++ {
		fmt.Printf("%s Please input your prompt: ", strconst.EmojiTips)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		userPrompt := scanner.Text()
		userPromptHistory = append(userPromptHistory, userPrompt)
		fullPrompt := strings.Join(userPromptHistory, strconst.NewLine)

		fmt.Printf("%s Thinking...\n", strconst.EmojiRunning)
		response := agent.SendMessage(fullPrompt)
		fmt.Println(response)
	}
}
