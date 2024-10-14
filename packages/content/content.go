package content

import "github.com/vinibgoulart/gitbook-postgresql-vectorize/packages/space"

type Content struct {
	ID      int
	Text    string
	SpaceId int
	Space   *space.Space `pg:"rel:has-one"`
}
