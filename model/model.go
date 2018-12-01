package model

import "time"

type Model struct {
	ID        int       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Response struct {
	Token   string      `json:"token"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   `json:"error"`
}
