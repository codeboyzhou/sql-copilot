package llm

type Model struct {
	Name          string
	BaseUrl       string
	ApiKeyEnvName string
}

var DeepSeekChat = Model{
	Name:          "deepseek-chat",
	BaseUrl:       "https://api.deepseek.com",
	ApiKeyEnvName: "DEEPSEEK_API_KEY",
}
