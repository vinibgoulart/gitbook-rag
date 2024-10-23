package session

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/vinibgoulart/gitbook-rag/packages/session"
)

func CreateNewSession(ctx *context.Context, db *bun.DB, w http.ResponseWriter) string {
	newSession := &session.Session{
		ID:    uuid.New().String(),
		Valid: true,
	}

	_, err := db.NewInsert().
		Model(newSession).Exec(*ctx)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return ""
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sessionId",
		Value:    newSession.ID,
		Path:     "/",
		MaxAge:   24 * 60 * 60, // 1 dia
		HttpOnly: true,
		Secure:   false,
	})

	return newSession.ID
}
