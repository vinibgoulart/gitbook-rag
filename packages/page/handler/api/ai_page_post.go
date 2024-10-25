package page

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/uptrace/bun"
	"github.com/vinibgoulart/gitbook-rag/packages/page"
	"github.com/vinibgoulart/gitbook-rag/packages/utils"
	"github.com/vinibgoulart/zius"
)

type AiPrompt struct {
	Prompt string `json:"prompt"`
}

type AiResponse struct {
	Response string `json:"response"`
}

func AiPagePost(ctx *context.Context, db *bun.DB) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		var aiPrompt AiPrompt

		errJsonDecode := utils.JsonDecode(res, req, &aiPrompt)
		if errJsonDecode != nil {
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		_, errValidate := zius.Validate(aiPrompt)
		if errValidate != nil {
			http.Error(res, errValidate.Error(), http.StatusBadRequest)
			return
		}

		response, err := page.GetResponseEmbeddingQuery(ctx, db)(&aiPrompt.Prompt)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")

		aiResponse := AiResponse{
			Response: response,
		}
		jsonResponse, err := json.Marshal(aiResponse)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Failed to encode response", http.StatusInternalServerError)
			return
		}

		res.Write(jsonResponse)
	}
}
