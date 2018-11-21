package model

import "net/http"

type Post struct {
	Model
	Name string `json:"name"`
	Body string `json:"body"`
	URL  string `json:"url"`
}

func NewPost(name, url string) *Post {
	return &Post{
		Name: name,
		URL:  url,
	}
}

func (p *Post) Bind(r *http.Request) error {
	return nil
}
