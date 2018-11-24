package model

import "time"

type Model struct {
	ID        int       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   `json:"error"`
}

func NewError(err error) *Error {
	return &Error{
		IsErr:   true,
		Message: err,
	}
}
