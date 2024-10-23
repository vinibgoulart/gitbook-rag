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
	ID        string           `bun:"id,pk" json:"id"`
	Agent     string           `bun:",notnull" json:"agent"`
	SessionId string           `bun:",notnull" json:"session_id"`
	Text      string           `bun:",notnull" json:"text"`
	Session   *session.Session `bun:"rel:has-one,join:session_id=id" json:"session"`
	CreatedAt time.Time        `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time        `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
}
