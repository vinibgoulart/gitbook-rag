package openai

import (
	"context"
	"fmt"

	_ "github.com/joho/godotenv/autoload"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/shared"
)

func GetEmbedding(text *string) []float64 {
	chatCompletion, err := Client.Embeddings.New(context.TODO(), openai.EmbeddingNewParams{
		Input: openai.F[openai.EmbeddingNewParamsInputUnion](shared.UnionString(*text)),
		Model: openai.F(openai.EmbeddingModelTextEmbeddingAda002),
	})

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return chatCompletion.Data[0].Embedding
}
