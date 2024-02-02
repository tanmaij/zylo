package openai

import "github.com/tanmaij/zylo/internal/model"

type chatCompletionRequest struct {
	Messages         []message    `json:"messages"`
	Model            string       `json:"model"`
	FrequencyPenalty *float64     `json:"frequency_penalty,omitempty"`
	LogitBias        *interface{} `json:"logit_bias,omitempty"` // You can adjust the type based on your needs
	Logprobs         *bool        `json:"logprobs,omitempty"`
	TopLogprobs      *interface{} `json:"top_logprobs,omitempty"` // You can adjust the type based on your needs
	MaxTokens        *int         `json:"max_tokens,omitempty"`
	N                *int         `json:"n,omitempty"`
	PresencePenalty  *float64     `json:"presence_penalty,omitempty"`
	ResponseFormat   *format      `json:"response_format,omitempty"`
	Seed             *interface{} `json:"seed,omitempty"` // You can adjust the type based on your needs
	Stop             *interface{} `json:"stop,omitempty"` // You can adjust the type based on your needs
	Stream           *bool        `json:"stream,omitempty"`
	Temperature      *float64     `json:"temperature,omitempty"`
	TopP             *float64     `json:"top_p,omitempty"`
	Tools            []tool       `json:"tools,omitempty"`
	User             *string      `json:"user,omitempty"`
}

type format struct {
	Type string `json:"type"`
}

type tool struct {
	ToolChoice string `json:"tool_choice"`
}

type chatCompletionResponse struct {
	ID      string      `json:"id"`
	Model   string      `json:"model"`
	Object  string      `json:"object"`
	Created int64       `json:"created"`
	Usage   usageDetail `json:"usage"`
	Choices []choice    `json:"choices"`
}

type usageDetail struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type choice struct {
	FinishReason string   `json:"finish_reason"`
	Index        int      `json:"index"`
	Message      message  `json:"message"`
	Logprobs     []string `json:"logprobs"`
}

type message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

func (m message) toModel() model.Message {
	return model.Message{
		Role:    m.Role,
		Content: m.Content,
	}
}
