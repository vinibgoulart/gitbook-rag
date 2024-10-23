package session

import "time"

type contextKey string

const SessionIDKey contextKey = "session_id"

type Session struct {
	ID        string `bun:"id,pk"`
	Context   string
	Valid     bool
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}
