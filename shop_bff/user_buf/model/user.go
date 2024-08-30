package model

import (
	"time"
)

type User struct {
	ID        int32     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
	Password  string    `json:"password"`
	NickName  string    `json:"nick_name"`
	Mobile    string    `json:"mobile"`
	Role      string    `json:"role"`
	Birthday  string    `json:"birthday"`
	Gender    string    `json:"gender"`
}
