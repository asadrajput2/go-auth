package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/asadrajput2/go-auth/pkg/models"
)

type Storage struct {
	*sql.DB
}

func Connect() (*Storage, error) {

	sqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5432, "asd", "whythis", "gotest")
	db, err := sql.Open("postgres", sqlInfo)

	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()
	return &Storage{db}, err
}

func (s *Storage) GetPosts(limit int, userId uint8) ([]models.Article, error) {
	var posts []models.Article
	err := s.QueryRow("SELECT * FROM posts WHERE author_id=$1 LIMIT $2", userId, limit).Scan(&posts)
	return posts, err
}

func (s *Storage) GetPost(id uint64) (models.Article, error) {

	var post models.Article
	err := s.QueryRow("SELECT * FROM posts WHERE id=$1", id).Scan(&post)
	return post, err
}

func (s *Storage) CreatePost(post models.Article) error {
	_, err := s.Exec("INSERT INTO posts (title, content, author_id) VALUES ($1, $2, $3)", post.Title, post.Conten, post.AuthorId)
	return err
}

func (s *Storage) UpdatePost(post models.Article) error {
	_, err := s.Exec("UPDATE posts SET title=$1, content=$2 WHERE id=$3", post.Title, post.Content, post.Id)
	return err
}

func (s *Storage) DeletePost(id uint64) error {
	_, err := s.Exec("DELETE FROM posts WHERE id=$1", id)
	return err
}

func (s *Storage) GetUsers() ([]models.User, error) {
	var users []models.User
	err := s.QueryRow("SELECT * FROM users").Scan(&users)
	return users, err
}

func (s *Storage) GetUser(id uint8) (models.User, error) {
	var user models.User
	err := s.QueryRow("SELECT * FROM users WHERE id=$1", id).Scan(&user)
	return user, err
}

func (s *Storage) AddUser(user models.User) (models.User, error) {
	err := s.QueryRow("INSERT INTO users (password, email, phone) VALUES ($1, $2, $3) RETURNING *", user.Email, user.Password, user.Phone).Scan(&user)
	return user, err
}
