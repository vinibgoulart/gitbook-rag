package content

import (
	"time"

	"github.com/vinibgoulart/gitbook-postgresql-vectorize/packages/space"
)

type Content struct {
	ID        string
	Title     string
	SpaceId   string
	Space     *space.Space `pg:"rel:has-one"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
