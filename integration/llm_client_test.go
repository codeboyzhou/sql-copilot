//go:build integration
// +build integration

package integration

import (
	"fmt"
	"testing"

	"github.com/joho/godotenv"

	"github.com/codeboyzhou/sql-copilot/agentic/llm"
	"github.com/codeboyzhou/sql-copilot/strconst"
)

func TestOpenAICompatibleClientSendMessage(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    string
	}{
		{
			name:    "TestSendMessage",
			message: "Hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := godotenv.Load("./testdata/.env"); err != nil {
				t.Errorf("No .env file found: %v", err)
			}

			var client llm.OpenAICompatibleClient
			client.InitializeWithModel(llm.DeepSeekChat)
			got := client.SendMessage(tt.message)
			fmt.Printf("SendMessage(%s), Received: %s\n", tt.message, got)

			if got == strconst.Empty {
				t.Errorf("SendMessage(%s) failed, got empty string, but want non-empty string", tt.message)
			}
		})
	}
}
