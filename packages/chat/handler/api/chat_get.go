package session

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/uptrace/bun"
	"github.com/vinibgoulart/gitbook-rag/packages/chat"
)

func ChatGet(ctx *context.Context, db *bun.DB) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		chatsWithSession, err := chat.GetSessionWithChats(ctx, db)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")

		response, err := json.Marshal(chatsWithSession)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Failed to encode response", http.StatusInternalServerError)
			return
		}

		res.Write(response)
	}
}
