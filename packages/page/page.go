package page

import (
	"time"

	"github.com/vinibgoulart/gitbook-postgresql-vectorize/packages/content"
	"github.com/vinibgoulart/gitbook-postgresql-vectorize/packages/space"
)

type Page struct {
	ID        string
	Text      string
	SpaceId   string
	Space     *space.Space `pg:"rel:has-one"`
	ContentId string
	Content   *content.Content `pg:"rel:has-one"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
