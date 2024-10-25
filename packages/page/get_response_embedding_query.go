package page

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/openai/openai-go"
	"github.com/pgvector/pgvector-go"
	"github.com/uptrace/bun"
	"github.com/vinibgoulart/gitbook-rag/packages/chat"
	openaiL "github.com/vinibgoulart/gitbook-rag/packages/openai"
	"github.com/vinibgoulart/gitbook-rag/packages/session"
	"github.com/vinibgoulart/gitbook-rag/packages/utils"
)

func GetResponseEmbeddingQuery(ctx *context.Context, db *bun.DB) func(query *string) (string, error) {
	return func(query *string) (string, error) {
		embed := openaiL.GetEmbedding(query)

		var items []Page
		err := db.NewSelect().
			Model(&items).
			OrderExpr("embedding <-> ?", pgvector.NewVector(utils.Float64ToFloat32(embed))).
			Limit(1).
			Scan(*ctx)
		if err != nil {
			return "", err
		}

		if len(items) == 0 {
			return "", fmt.Errorf("no page found")
		}

		ctxSessionId, ok := (*ctx).Value(session.SessionIDKey).(string)
		if !ok {
			return "", fmt.Errorf("session_id not found or not a string")
		}

		var session session.Session
		errSession := db.NewSelect().
			Model(&session).
			Where("id = ?", ctxSessionId).
			Scan(*ctx)
		if errSession != nil {
			return "", errSession
		}

		if session.Context == "" {
			session.Context = items[0].Text
			db.NewUpdate().
				Model(&session).
				Where("id = ?", ctxSessionId).
				Exec(*ctx)
		}

		_, errChatsInsert := db.NewInsert().
			Model(&chat.Chat{
				ID:        uuid.New().String(),
				SessionId: ctxSessionId,
				Agent:     "user",
				Text:      *query,
			}).
			Exec(*ctx)
		if errChatsInsert != nil {
			return "", errChatsInsert
		}

		var chats []chat.Chat
		errChats := db.NewSelect().
			Model(&chats).
			Where("session_id = ?", ctxSessionId).
			Scan(*ctx)
		if errChats != nil {
			return "", errChats
		}

		var allMessages []openai.ChatCompletionMessageParamUnion
		if len(chats) == 0 {
			for _, chat := range chats {
				if chat.Agent == "user" {
					allMessages = append(allMessages, openai.UserMessage(chat.Text))
				}
				if chat.Agent == "assistant" {
					allMessages = append(allMessages, openai.AssistantMessage(chat.Text))
				}
			}
		}

		res, err := openaiL.GenerateCompletion(ctx)(&session.Context, allMessages...)
		if err != nil {
			return "", err
		}

		_, errChatsInsert = db.NewInsert().
			Model(&chat.Chat{
				ID:        uuid.New().String(),
				SessionId: ctxSessionId,
				Agent:     "assistant",
				Text:      res,
			}).
			Exec(*ctx)
		if errChatsInsert != nil {
			return "", errChatsInsert
		}

		return res, nil
	}
}
