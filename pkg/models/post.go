package models

type Article struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	AuthorId uint8  `json:"author_id"`
}
