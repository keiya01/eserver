package model

type Post struct {
	Model
	Title string `json:"title"`
	Body  string `json:"body"`
}

type PostService interface {
	FindAll() (*Post, error)
}

func NewPost(title, body string) *Post {
	return &Post{
		Title: title,
		Body:  body,
	}
}
