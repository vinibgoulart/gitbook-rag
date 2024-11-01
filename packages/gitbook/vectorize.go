package gitbook

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pgvector/pgvector-go"
	"github.com/uptrace/bun"
	"github.com/vinibgoulart/gitbook-rag/packages/content"
	"github.com/vinibgoulart/gitbook-rag/packages/database"
	"github.com/vinibgoulart/gitbook-rag/packages/openai"
	"github.com/vinibgoulart/gitbook-rag/packages/page"
	"github.com/vinibgoulart/gitbook-rag/packages/space"
	"github.com/vinibgoulart/gitbook-rag/packages/utils"
)

func Vectorize(ctx *context.Context, db *bun.DB) error {
	spaces, errSpacesGet := SpacesGet()

	if errSpacesGet != nil {
		return errSpacesGet
	}

	var items []SpaceItems

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
	} else {
		items = spaces.Items
	}

	for _, s := range items {
		errSpaceCreate := database.InsertOrUpdate(ctx, db, &space.Space{
			ID:    s.ID,
			Title: s.Title,
			Url:   s.Urls.Published,
		}, "id", "title = EXCLUDED.title, updated_at = NOW(), url = EXCLUDED.url")

		if errSpaceCreate != nil {
			fmt.Println(errSpaceCreate.Error())
			return errSpaceCreate
		}

		c, errSpaceContentGet := ContentGet(&s.ID)

		if errSpaceContentGet != nil {
			return errSpaceContentGet
		}

		var allIdsUpdate []string
		for _, p := range c.Pages {
			allIdsUpdate = append(allIdsUpdate, utils.RecursiveCatchFields("id", p)...)
		}

		for _, id := range allIdsUpdate {
			errContentCreate := database.InsertOrUpdate(ctx, db, &content.Content{
				ID:      id,
				SpaceId: s.ID,
			}, "id", "space_id = EXCLUDED.space_id, updated_at = NOW()")

			if errContentCreate != nil {
				fmt.Println(errContentCreate.Error())
				return errContentCreate
			}

			allIdsUpdate = append(allIdsUpdate, id)

			go VectorizePages(ctx, db)(&s.Urls.Published, &s.ID, &id)
		}

		var pagesNotUpdated []page.Page
		errSelect := db.NewSelect().Model(&pagesNotUpdated).Where("space_id = ?", s.ID).Where("content_id NOT IN (?)", allIdsUpdate).Scan(*ctx)
		if errSelect != nil {
			fmt.Println(errSelect.Error())
			return errSelect
		}

		for _, p := range pagesNotUpdated {
			_, errDelete := db.NewDelete().Model(&p).Where("id = ?", p.ID).Exec(*ctx)
			if errDelete != nil {
				fmt.Println(errDelete.Error())
				return errDelete
			}
		}
	}

	return nil
}

func VectorizePages(ctx *context.Context, db *bun.DB) func(spaceUrl *string, spaceId *string, contentPageId *string) error {
	return func(spaceUrl *string, spaceId *string, contentPageId *string) error {
		p, errPageContentGet := SpacesContentPageGet(spaceId, contentPageId)

		if errPageContentGet != nil {
			fmt.Println(errPageContentGet.Error())
			return errPageContentGet
		}

		if p.Type != "document" {
			return nil
		}

		var text string
		if p.Description != "" {
			text = p.Title + ". " + p.Description
		} else {
			text = p.Title + ". "
		}

		var allTexts []string
		for _, node := range p.Document.Nodes {
			allTexts = append(allTexts, utils.RecursiveCatchFields("text", node)...)
		}

		if len(allTexts) < 1 {
			return nil
		}

		text = text + " " + strings.Join(allTexts, " ")

		if text != "" {
			embed := openai.GetEmbedding(&text)
			errPageCreate := database.InsertOrUpdate(ctx, db, &page.Page{
				ID:        p.ID,
				Text:      text,
				Title:     p.Title,
				Url:       fmt.Sprintf("%s%s", *spaceUrl, p.Path),
				SpaceId:   *spaceId,
				ContentId: *contentPageId,
				Embedding: pgvector.NewVector(utils.Float64ToFloat32(embed)),
			}, "id", "text = EXCLUDED.text, space_id = EXCLUDED.space_id, content_id = EXCLUDED.content_id, embedding = EXCLUDED.embedding, updated_at = NOW(), url = EXCLUDED.url, title = EXCLUDED.title")

			if errPageCreate != nil {
				fmt.Println(errPageCreate.Error())
				return errPageCreate
			}
		}

		return nil
	}
}
