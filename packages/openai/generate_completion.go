package openai

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go"
)

func GenerateCompletion(ctx *context.Context) func(context *string, messages ...openai.ChatCompletionMessageParamUnion) (string, error) {
	return func(context *string, messages ...openai.ChatCompletionMessageParamUnion) (string, error) {
		systemMessage := openai.SystemMessage(fmt.Sprintf("You are a helpful assistant that uses the provided context as reference, but can answer questions that are outside the scope of the context in a natural way. Always use the %s language. Context: %s", os.Getenv("CHATBOT_LANGUAGE"), *context))

		allMessages := append([]openai.ChatCompletionMessageParamUnion{
			systemMessage,
		}, messages...)

		chatCompletion, err := Client.Chat.Completions.New(*ctx, openai.ChatCompletionNewParams{
			Messages:    openai.F(allMessages),
			Model:       openai.F(openai.ChatModelGPT4oMini),
			Temperature: openai.Float(0.5),
		})

		if err != nil {
			return "", err
		}

		return chatCompletion.Choices[0].Message.Content, nil
	}
}
