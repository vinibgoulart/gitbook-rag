package gitbook

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-pg/pg"
	"github.com/vinibgoulart/gitbook-postgresql-vectorize/packages/content"
	"github.com/vinibgoulart/gitbook-postgresql-vectorize/packages/database"
	"github.com/vinibgoulart/gitbook-postgresql-vectorize/packages/page"
	"github.com/vinibgoulart/gitbook-postgresql-vectorize/packages/space"
	"github.com/vinibgoulart/gitbook-postgresql-vectorize/packages/utils"
)

func Vectorize(db *pg.DB) error {
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
		errSpaceCreate := database.InsertOrUpdate(db, &space.Space{
			ID:    s.ID,
			Title: s.Title,
		}, "id", "title = EXCLUDED.title")

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

				errContentCreate := database.InsertOrUpdate(db, &content.Content{
					ID:      cp.ID,
					Title:   cp.Title,
					SpaceId: s.ID,
				}, "id", "title = EXCLUDED.title, space_id = EXCLUDED.space_id")

				if errContentCreate != nil {
					fmt.Println(errContentCreate.Error())
					return errContentCreate
				}

				go VectorizePages(db)(&s.ID, &cp.ID)
			}
		}
	}

	return nil
}

func VectorizePages(db *pg.DB) func(spaceId *string, contentPageId *string) error {
	return func(spaceId *string, contentPageId *string) error {
		p, errPageContentGet := SpacesContentPageGet(spaceId, contentPageId)

		if errPageContentGet != nil {
			fmt.Println(errPageContentGet.Error())
			return errPageContentGet
		}

		var text string

		for _, node := range p.Document.Nodes {
			textCurrent := utils.RecursiveCatchField("text", node)
			if textCurrent != "" {
				text = text + " " + textCurrent
			}
		}

		if text != "" {
			errPageCreate := database.InsertOrUpdate(db, &page.Page{
				ID:        p.ID,
				Text:      text,
				SpaceId:   *spaceId,
				ContentId: *contentPageId,
			}, "id", "text = EXCLUDED.text, space_id = EXCLUDED.space_id, content_id = EXCLUDED.content_id")

			if errPageCreate != nil {
				fmt.Println(errPageCreate.Error())
				return errPageCreate
			}
		}

		return nil
	}
}
