package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/vinibgoulart/gitbook-rag/packages/session"
)

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
				sessionId = createNewSession(ctx, db, w)
			} else {
				sessionId = cookie.Value
				var sess session.Session
				err = db.NewSelect().
					Model(&sess).
					Where("id = ?", sessionId).
					Scan(*ctx)

				if err != nil || !sess.Valid {
					if err == nil {
						sess.Valid = false
						_, err = db.NewUpdate().
							Model(&sess).
							Where("id = ?", sessionId).
							Exec(*ctx)
						if err != nil {
							fmt.Println("Error")
						}
					}

					sessionId = createNewSession(ctx, db, w)
				}
			}

			*ctx = context.WithValue(*ctx, session.SessionIDKey, sessionId)
			next.ServeHTTP(w, r)
		})
	}
}

func createNewSession(ctx *context.Context, db *bun.DB, w http.ResponseWriter) string {
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
