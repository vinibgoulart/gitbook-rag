package session

import (
	"context"
	"fmt"
	"net/http"

	"github.com/uptrace/bun"
	"github.com/vinibgoulart/gitbook-rag/packages/session"
)

func RemoveSession(ctx *context.Context, db *bun.DB, res http.ResponseWriter) {
	ctxSessionId, ok := (*ctx).Value(session.SessionIDKey).(string)
	if !ok {
		fmt.Println("Error")
	}

	*ctx = context.WithValue(*ctx, session.SessionIDKey, "")

	http.SetCookie(res, &http.Cookie{
		Name:     "sessionId",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	})

	session := session.Session{
		Valid: false,
	}

	_, err := db.NewUpdate().
		Model(&session).
		Column("valid").
		Where("id = ?", ctxSessionId).
		Exec(*ctx)

	if err != nil {
		fmt.Println(err)
	}
}
