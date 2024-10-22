package page

import (
	"time"

	"github.com/vinibgoulart/gitbook-postgresql-vectorize/packages/content"
	"github.com/vinibgoulart/gitbook-postgresql-vectorize/packages/space"
)

type Page struct {
	ID        string `bun:"id,pk"`
	Text      string `bun:"text,notnull"`
	SpaceId   string
	Space     *space.Space `bun:"rel:has-one,join:space_id=id"`
	ContentId string
	Content   *content.Content `bun:"rel:has-one,join:content_id=id"`
	CreatedAt time.Time        `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time        `bun:",nullzero,notnull,default:current_timestamp"`
}
