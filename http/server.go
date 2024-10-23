package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/uptrace/bun"
	httpSession "github.com/vinibgoulart/gitbook-rag/http/session"
	chat "github.com/vinibgoulart/gitbook-rag/packages/chat/handler/api"
	page "github.com/vinibgoulart/gitbook-rag/packages/page/handler/api"
)

func ServerInit(db *bun.DB) func(ctx context.Context, waitGroup *sync.WaitGroup) {
	return func(ctx context.Context, waitGroup *sync.WaitGroup) {
		defer waitGroup.Done()

		router := chi.NewRouter()

		router.Use(middleware.RequestID)
		router.Use(middleware.RealIP)
		router.Use(middleware.Logger)
		router.Use(middleware.Recoverer)
		router.Use(JsonContentTypeMiddleware)
		router.Use(SessionMiddleware(&ctx, db))
		router.Use(middleware.Timeout(30 * time.Second))

		router.Get("/status", func(res http.ResponseWriter, req *http.Request) {
			res.Write([]byte("OK"))
		})
		router.Post("/logout", func(res http.ResponseWriter, req *http.Request) {
			httpSession.RemoveSession(&ctx, db, res)
		})
		router.Post("/ai/page", page.AiPagePost(&ctx, db))
		router.Get("/chat", chat.ChatGet(&ctx, db))

		port := fmt.Sprintf(":%s", os.Getenv("PORT"))

		server := &http.Server{
			Addr:    port,
			Handler: router,
		}

		go func() {
			fmt.Println("HTTP server listening on", port)
			server.ListenAndServe()
		}()

		<-ctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatalf("HTTP server shutdown error: %s", err)
		}
	}
}
