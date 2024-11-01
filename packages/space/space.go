package space

import "time"

type Space struct {
	ID        string    `bun:"id,pk"`
	Title     string    `bun:"title,notnull"`
	Url       string    `bun:"url"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}
