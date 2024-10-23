package session

import (
	"context"
	"net/http"

	"github.com/uptrace/bun"
	"github.com/vinibgoulart/gitbook-rag/packages/session"
)

func RemoveSession(ctx *context.Context, db *bun.DB, res http.ResponseWriter) {
	*ctx = context.WithValue(*ctx, session.SessionIDKey, "")

	http.SetCookie(res, &http.Cookie{
		Name:     "sessionId",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	})
}
