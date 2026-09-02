package entity

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:product"`

	ID           int64     `bun:"id,pk,autoincrement"`
	GUID         uuid.UUID `bun:"guid,pk,notnull"`
	CategoryGUID uuid.UUID `bun:"category_guid,notnull"`
	Name         string    `bun:"name,notnull"`
	Description  *string   `bun:"description"`              // указатель, может быть NULL
	Price        *float64  `bun:"price,type:numeric(10,2)"` // указатель
	CreatedAt    time.Time `bun:"created_at,notnull,default:now()"`
	UpdatedAt    time.Time `bun:"updated_at,notnull,default:now()"`
}
