package gitbook

import (
	"fmt"
	"os"
	"strings"

	"github.com/vinibgoulart/gitbook-postgresql-vectorize/packages/utils"
)

func Vectorize() error {
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

	for _, space := range items {
		spaceContent, errSpaceContentGet := SpacesContentGet(&space.ID)

		if errSpaceContentGet != nil {
			return errSpaceContentGet
		}

		for _, page := range spaceContent.Pages {
			for _, page := range page.Pages {
				go VectorizePages(&space.ID, &page.ID)
			}
		}
	}

	return nil
}

func VectorizePages(spaceId *string, pageId *string) error {
	pageContent, errPageContentGet := SpacesContentPageGet(spaceId, pageId)

	if errPageContentGet != nil {
		fmt.Println(errPageContentGet.Error())
		return errPageContentGet
	}

	for _, node := range pageContent.Document.Nodes {
		text := utils.RecursiveCatchField("text", node)

		if textStr, ok := text.(string); ok && textStr != "" {
			fmt.Println(text)
		}
	}

	return nil
}
