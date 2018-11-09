package model

type Post struct {
	Model
	Title string `json:"title"`
	Body  string `json:"body"`
}

func NewPost(title, body string) *Post {
	return &Post{
		Title: title,
		Body:  body,
	}
}
