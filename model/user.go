package model

type User struct {
	Model
	Email    string `json:"email" gorm:"unique_index"`
	Password string `json:"password"`
}
