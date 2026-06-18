package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Experiment struct {
	bun.BaseModel `bun:"table:experiments,alias:exp"`

	ID         int64   `bun:"id,pk,autoincrement" form:"id" json:"id" xml:"id"`
	SubjectID  int64   `bun:"subject_id" form:"subject_id" json:"subject_id" xml:"subject_id"`
	AdminIDS   []int64 `bun:"admin_ids,array" form:"admin_ids" json:"admin_ids" xml:"admin_ids"`
	TeacherIDS []int64 `bun:"teacher_ids,array" form:"teacher_ids" json:"teacher_ids" xml:"teacher_ids"`
	StudentIDS []int64 `bun:"student_ids,array" form:"student_ids" json:"student_ids" xml:"student_ids"`

	Image       string `bun:"image" form:"image" json:"image" xml:"image"`
	Title       string `bun:"title" form:"title" json:"title" xml:"title"`
	Description string `bun:"description" form:"description" json:"description" xml:"description"`
	Status      string `bun:"status" form:"status" json:"status" xml:"status"`

	Rules  []string `bun:"rules,array" form:"rules" json:"rules" xml:"rules"`
	Groups []any    `bun:"groups,type:jsonb" form:"groups" json:"groups" xml:"groups"`

	Questions   []string       `bun:"questions,array" form:"questions" json:"questions" xml:"questions"`
	Hypotesis   []string       `bun:"hypotesis,array" form:"hypotesis" json:"hypotesis" xml:"hypotesis"`
	Simulation  map[string]any `bun:"simulation,type:jsonb" form:"simulation" json:"simulation" xml:"simulation"`
	Data        map[string]any `bun:"data,type:jsonb" form:"data" json:"data" xml:"data"`
	Explanation map[string]any `bun:"explanation,type:jsonb" form:"explanation" json:"explanation" xml:"explanation"`

	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" form:"created_at" json:"created_at" xml:"created_at"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" form:"updated_at" json:"updated_at" xml:"updated_at"`
}
