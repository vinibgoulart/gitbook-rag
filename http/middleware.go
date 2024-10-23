package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/vinibgoulart/gitbook-rag/packages/session"
)

type contextKey string

const sessionIDKey contextKey = "sessionId"

func JsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
func SessionMiddleware(ctx *context.Context, db *bun.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("sessionId")
			var sessionId string

			if err != nil || cookie.Value == "" {
				newSession := &session.Session{
					ID: uuid.New().String(),
				}

				_, err := db.NewInsert().
					Model(newSession).Exec(*ctx)

				if err != nil {
					fmt.Println(err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				sessionId = newSession.ID

				http.SetCookie(w, &http.Cookie{
					Name:     "sessionId",
					Value:    sessionId,
					Path:     "/",
					MaxAge:   24 * 60 * 60, // 1 day
					HttpOnly: true,
					Secure:   false,
				})
				*ctx = context.WithValue(*ctx, sessionIDKey, sessionId)
			}

			next.ServeHTTP(w, r)
		})
	}
}
