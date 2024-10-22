package page

import (
	"context"

	"github.com/pgvector/pgvector-go"
	"github.com/uptrace/bun"
	"github.com/vinibgoulart/gitbook-rag/packages/openai"
	"github.com/vinibgoulart/gitbook-rag/packages/utils"
)

func GetResponseEmbeddingQuery(ctx *context.Context, db *bun.DB) func(query *string) (string, error) {
	return func(query *string) (string, error) {
		embed := openai.GetEmbedding(query)

		var items []Page
		err := db.NewSelect().
			Model(&items).
			OrderExpr("embedding <-> ?", pgvector.NewVector(utils.Float64ToFloat32(embed))).
			Limit(1).
			Scan(*ctx)

		if err != nil {
			return "", err
		}

		res, err := openai.GenerateCompletion(ctx)(&items[0].Text, query)
		if err != nil {
			return "", err
		}

		return res, nil
	}
}
