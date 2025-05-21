package user

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	Id        string           `json:"id"`
	Username  string           `json:"username"`
	Email     string           `json:"email"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	DeletedAt pgtype.Timestamp `json:"deleted_at"`
	Deleted   bool             `json:"deleted"`
}
