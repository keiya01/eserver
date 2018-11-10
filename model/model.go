package model

import "time"

type Model struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Error struct {
	IsErr   bool  `json:"isErr"`
	Message error `json:"message"`
}

func NewError(err error) *Error {
	return &Error{
		IsErr:   true,
		Message: err,
	}
}
