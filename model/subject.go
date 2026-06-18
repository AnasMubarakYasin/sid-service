package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Subject struct {
	bun.BaseModel `bun:"table:subjects,alias:sbj"`

	ID          int64    `bun:"id,pk,autoincrement" form:"id" json:"id" xml:"id"`
	Title       string   `bun:"title,notnull" form:"title" json:"title" xml:"title"`
	Description string   `bun:"description,notnull" form:"description" json:"description" xml:"description"`
	Video       string   `bun:"video,notnull" form:"video" json:"video" xml:"video"`
	Questions   []string `bun:"questions,notnull" form:"questions" json:"questions" xml:"questions"`

	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" form:"created_at" json:"created_at" xml:"created_at"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" form:"updated_at" json:"updated_at" xml:"updated_at"`
}
