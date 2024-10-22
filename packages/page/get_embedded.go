package page

import (
	"context"

	"github.com/pgvector/pgvector-go"
	"github.com/uptrace/bun"
	"github.com/vinibgoulart/gitbook-llm/packages/openai"
	"github.com/vinibgoulart/gitbook-llm/packages/utils"
)

func GetEmbedded(ctx *context.Context, db *bun.DB) func(query *string) Page {
	return func(query *string) Page {
		embed := openai.GetEmbedding(query)

		var items []Page
		err := db.NewSelect().
			Model(&items).
			OrderExpr("embedding <-> ?", pgvector.NewVector(utils.Float64ToFloat32(embed))).
			Limit(1).
			Scan(*ctx)

		if err != nil {
			panic(err)
		}

		return items[0]
	}
}
