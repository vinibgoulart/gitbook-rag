package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/uptrace/bun"
	httpSession "github.com/vinibgoulart/gitbook-rag/http/session"
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
				sessionId = httpSession.CreateNewSession(ctx, db, w)
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

					sessionId = httpSession.CreateNewSession(ctx, db, w)
				}
			}

			*ctx = context.WithValue(*ctx, session.SessionIDKey, sessionId)
			next.ServeHTTP(w, r)
		})
	}
}
