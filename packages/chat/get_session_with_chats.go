package chat

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/vinibgoulart/gitbook-rag/packages/session"
)

type SessionWithChats struct {
	Context   string `json:"context"`
	CreatedAt string `json:"created_at"`
	Chats     []Chat `json:"chats"`
}

func GetSessionWithChats(ctx *context.Context, db *bun.DB) (*SessionWithChats, error) {
	ctxSessionId, ok := (*ctx).Value(session.SessionIDKey).(string)
	if !ok {
		return nil, fmt.Errorf("session_id not found or not a string")
	}

	var session session.Session
	errSession := db.NewSelect().
		Model(&session).
		Where("id = ?", ctxSessionId).
		Scan(*ctx)
	if errSession != nil {
		return nil, errSession
	}

	var chats []Chat
	errChats := db.NewSelect().
		Model(&chats).
		Where("session_id = ?", ctxSessionId).
		Scan(*ctx)
	if errChats != nil {
		return nil, errChats
	}

	return &SessionWithChats{
		Context:   session.Context,
		CreatedAt: session.CreatedAt.String(),
		Chats:     chats,
	}, nil
}
