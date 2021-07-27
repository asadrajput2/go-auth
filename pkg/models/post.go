package models

type Article struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"desc"`
	Content     string `json:"content"`
	AuthorId    string `json:"author_id"`
}
