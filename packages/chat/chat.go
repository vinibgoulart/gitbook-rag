package chat

import (
	"time"

	"github.com/vinibgoulart/gitbook-rag/packages/session"
)

const (
	AgentTypeUser = "user"
	AgentTypeBot  = "bot"
)

type Chat struct {
	ID        string `bun:"id,pk"`
	Agent     string `bun:",notnull"`
	SessionId string
	Text      string
	Session   *session.Session `bun:"rel:has-one,join:session_id=id"`
	CreatedAt time.Time        `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time        `bun:",nullzero,notnull,default:current_timestamp"`
}
