package page

import (
	"context"

	"github.com/pgvector/pgvector-go"
	"github.com/uptrace/bun"
	"github.com/vinibgoulart/gitbook-llm/packages/openai"
	"github.com/vinibgoulart/gitbook-llm/packages/utils"
)

func GetEmbedded(ctx *context.Context, db *bun.DB) func(query *string) (Page, error) {
	return func(query *string) (Page, error) {
		embed := openai.GetEmbedding(query)

		var items []Page
		err := db.NewSelect().
			Model(&items).
			OrderExpr("embedding <-> ?", pgvector.NewVector(utils.Float64ToFloat32(embed))).
			Limit(1).
			Scan(*ctx)

		if err != nil {
			return Page{}, err
		}

		return items[0], nil
	}
}
