package openai

import (
	"context"

	"github.com/openai/openai-go"
)

func GenerateCompletion(ctx *context.Context) func(context *string, messages ...openai.ChatCompletionMessageParamUnion) (string, error) {
	return func(context *string, messages ...openai.ChatCompletionMessageParamUnion) (string, error) {
		allMessages := append([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(*context),
		}, messages...)

		chatCompletion, err := Client.Chat.Completions.New(*ctx, openai.ChatCompletionNewParams{
			Messages:    openai.F(allMessages),
			Model:       openai.F(openai.ChatModelGPT4o),
			Temperature: openai.Float(0.3),
		})

		if err != nil {
			return "", err
		}

		return chatCompletion.Choices[0].Message.Content, nil
	}
}
