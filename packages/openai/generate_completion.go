package openai

import (
	"context"

	"github.com/openai/openai-go"
)

func GenerateCompletion(ctx *context.Context) func(prompt *string, context *string) (string, error) {
	return func(prompt *string, context *string) (string, error) {
		chatCompletion, err := Client.Chat.Completions.New(*ctx, openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(*context),
				openai.UserMessage(*prompt),
			}),
			Model:       openai.F(openai.ChatModelGPT4o),
			Temperature: openai.Float(0.3),
		})

		if err != nil {
			return "", err
		}

		return chatCompletion.Choices[0].Message.Content, nil
	}
}
