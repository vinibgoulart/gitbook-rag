package openai

import (
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var (
	Client = openai.NewClient(
		option.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
	)
)
