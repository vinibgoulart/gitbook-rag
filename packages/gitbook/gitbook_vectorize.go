package gitbook

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pgvector/pgvector-go"
	"github.com/uptrace/bun"
	"github.com/vinibgoulart/gitbook-llm/packages/content"
	"github.com/vinibgoulart/gitbook-llm/packages/database"
	"github.com/vinibgoulart/gitbook-llm/packages/openai"
	"github.com/vinibgoulart/gitbook-llm/packages/page"
	"github.com/vinibgoulart/gitbook-llm/packages/space"
	"github.com/vinibgoulart/gitbook-llm/packages/utils"
)

func Vectorize(ctx *context.Context, db *bun.DB) error {
	spaces, errSpacesGet := SpacesGet()

	if errSpacesGet != nil {
		return errSpacesGet
	}

	var items []GitbookSpaceItems

	if os.Getenv("GITBOOK_SPACES_IDS") != "" {
		spacesIds := os.Getenv("GITBOOK_SPACES_IDS")

		spacesIdsArray := strings.Split(spacesIds, " ")
		for _, spaceId := range spacesIdsArray {
			for _, space := range spaces.Items {
				if space.ID == spaceId {
					items = append(items, space)
				}
			}
		}
	}

	for _, s := range items {
		errSpaceCreate := database.InsertOrUpdate(ctx, db, &space.Space{
			ID:    s.ID,
			Title: s.Title,
		}, "id", "title = EXCLUDED.title, updated_at = NOW()")

		if errSpaceCreate != nil {
			fmt.Println(errSpaceCreate.Error())
			return errSpaceCreate
		}

		c, errSpaceContentGet := SpacesContentGet(&s.ID)

		if errSpaceContentGet != nil {
			return errSpaceContentGet
		}

		for _, page := range c.Pages {
			for _, cp := range page.Pages {

				errContentCreate := database.InsertOrUpdate(ctx, db, &content.Content{
					ID:      cp.ID,
					Title:   cp.Title,
					SpaceId: s.ID,
				}, "id", "title = EXCLUDED.title, space_id = EXCLUDED.space_id, updated_at = NOW()")

				if errContentCreate != nil {
					fmt.Println(errContentCreate.Error())
					return errContentCreate
				}

				go VectorizePages(ctx, db)(&s.ID, &cp.ID)
			}
		}
	}

	return nil
}

func VectorizePages(ctx *context.Context, db *bun.DB) func(spaceId *string, contentPageId *string) error {
	return func(spaceId *string, contentPageId *string) error {
		p, errPageContentGet := SpacesContentPageGet(spaceId, contentPageId)

		if errPageContentGet != nil {
			fmt.Println(errPageContentGet.Error())
			return errPageContentGet
		}

		var text string
		if p.Description != "" {
			text = p.Title + ". " + p.Description
		} else {
			text = p.Title + ". "
		}

		for _, node := range p.Document.Nodes {
			textCurrent := utils.RecursiveCatchField("text", node)
			if textCurrent != "" {
				text = text + " " + textCurrent
			}

			if textCurrent != "" {
				embed := openai.GetEmbedding(&text)
				errPageCreate := database.InsertOrUpdate(ctx, db, &page.Page{
					ID:        p.ID,
					Text:      text,
					SpaceId:   *spaceId,
					ContentId: *contentPageId,
					Embedding: pgvector.NewVector(utils.Float64ToFloat32(embed)),
				}, "id", "text = EXCLUDED.text, space_id = EXCLUDED.space_id, content_id = EXCLUDED.content_id, embedding = EXCLUDED.embedding, updated_at = NOW()")

				if errPageCreate != nil {
					fmt.Println(errPageCreate.Error())
					return errPageCreate
				}
			}
		}

		return nil
	}
}
