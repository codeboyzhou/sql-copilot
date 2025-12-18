package llm

import (
	"context"
	"log"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"

	"github.com/codeboyzhou/sql-copilot/strconst"
)

type OpenAICompatibleClient struct {
	model  Model
	client openai.Client
}

func (c *OpenAICompatibleClient) InitializeWithModel(model Model) {
	apiKey := os.Getenv(model.ApiKeyEnvName)
	if apiKey == strconst.Empty {
		log.Fatalf("%s is not set", model.ApiKeyEnvName)
	}

	c.model = model
	c.client = openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(model.BaseUrl),
	)
}

func (c *OpenAICompatibleClient) SendMessage(message string) string {
	if c.model.Name == strconst.Empty {
		log.Fatalln("Client is not initialized")
	}

	ctx := context.TODO()
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(SystemPrompt),
		openai.UserMessage(message),
	}
	body := openai.ChatCompletionNewParams{
		Model:    c.model.Name,
		Messages: messages,
	}
	chatCompletion, err := c.client.Chat.Completions.New(ctx, body)
	if err != nil {
		log.Fatalln(err)
	}
	return chatCompletion.Choices[0].Message.Content
}
