package model

type Post struct {
	Model
	Name string `json:"name"`
	URL  string `json:"url"`
}

func NewPost(name, url string) *Post {
	return &Post{
		Name: name,
		URL:  url,
	}
}
