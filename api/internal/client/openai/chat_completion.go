package openai

import (
	"context"
	"log"

	"github.com/tanmaij/zylo/config"
	"github.com/tanmaij/zylo/internal/model"
	"github.com/tanmaij/zylo/pkg/api"
	"github.com/tanmaij/zylo/pkg/utils"
)

type ChatCompletion struct {
	apiKey               string
	model                string
	textGenerationClient *api.Client
}

func NewChatCompletion(apiKey, model string) (ChatCompletion, error) {
	textGenerationClient, err := api.NewClient(
		config.Instance.OpenAI.ChatCompletionAPIURL,
		api.MethodPost,
		"Generating text from chat completions",
		api.ContentTypeApplicationJSON)
	if err != nil {
		return ChatCompletion{}, err
	}

	return ChatCompletion{
		model:                model,
		apiKey:               apiKey,
		textGenerationClient: textGenerationClient,
	}, nil
}

func createPayload(body []byte) api.Payload {
	log.Printf("creating payload with body %v", string(body))
	return api.Payload{
		Body:        body,
		QueryParams: nil,
		PathVars:    nil,
		Header: map[string]string{
			"Authorization": "Bearer " + config.Instance.OpenAI.APIKey,
		},
	}
}

type ChatCompletionInput struct {
	Messages []model.Message
}

func (i ChatCompletionInput) toRequestBody() chatCompletionRequest {
	var msgs = make([]message, len(i.Messages))
	for idx := range i.Messages {
		msgs[idx] = message{
			Content: i.Messages[idx].Content,
			Role:    i.Messages[idx].Role,
		}
	}

	return chatCompletionRequest{
		Messages: msgs,
	}
}

func (c ChatCompletion) RequestToGenerate(ctx context.Context, input ChatCompletionInput) ([]model.Message, error) {
	requestBody := input.toRequestBody()
	requestBody.Model = c.model

	encodedRequestBody, err := utils.AnyToJSON(requestBody)
	if err != nil {
		log.Printf("ChatCompletion.RequestToGenerate: Error when encoding input: %v", err)
		return nil, err
	}

	resp, err := c.textGenerationClient.Send(ctx, createPayload(encodedRequestBody))
	if err != nil {
		log.Printf("ChatCompletion.RequestToGenerate: Error when sending request: %v", err)
		return nil, err
	}

	if resp.Status != api.StatusOK {
		log.Printf("ChatCompletion.RequestToGenerate: unexpected status code when sending request: %v", err)
		return nil, ErrUnexepectedCode
	}

	var decodedRespBody chatCompletionResponse
	if err := utils.JSONToAny(resp.Body, &decodedRespBody); err != nil {
		log.Printf("ChatCompletion.RequestToGenerate: Error when decoding response body: %v", err)
		return nil, err
	}

	var rs []model.Message
	for i := range decodedRespBody.Choices {
		switch decodedRespBody.Choices[i].FinishReason {
		case "stop": // good result
			log.Printf("good result")
			rs = append(rs, decodedRespBody.Choices[i].Message.toModel())
		case "length": // maximum number of tokens
			log.Printf("maximum number of tokens")
			rs = append(rs, decodedRespBody.Choices[i].Message.toModel())
		case "content_filter": // content filter
			log.Printf("content filter")
			rs = append(rs, decodedRespBody.Choices[i].Message.toModel())
		}
	}

	return rs, nil
}
