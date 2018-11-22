package model

type Post struct {
	Model
	Name string `json:"name"`
	Body string `json:"body"`
	URL  string `json:"url"`
}
