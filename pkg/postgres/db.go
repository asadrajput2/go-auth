package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/asadrajput2/go-auth/pkg/models"
	_ "github.com/lib/pq"
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

func (s *Storage) GetPosts(limit int, userId interface{}) ([]models.Article, error) {
	var posts []models.Article
	rows, err := s.Query("SELECT * FROM posts WHERE author_id=$1", userId)

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var article models.Article
		err := rows.Scan(&article.Id, &article.Title, &article.Content, &article.AuthorId)
		if err != nil {
			log.Fatal(err)
		}
		posts = append(posts, article)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return posts, err
}

func (s *Storage) GetPost(id interface{}) (models.Article, error) {

	var post models.Article
	err := s.QueryRow("SELECT * FROM posts WHERE id=$1", id).Scan(&post.Id, &post.Title, &post.Content, &post.AuthorId)
	return post, err
}

func (s *Storage) CreatePost(post models.Article, userId interface{}) error {
	_, err := s.Exec("INSERT INTO posts (title, content, author_id) VALUES ($1, $2, $3)", post.Title, post.Content, userId)
	return err
}

func (s *Storage) UpdatePost(post models.Article, postId string) error {
	_, err := s.Exec("UPDATE posts SET title=$1, content=$2 WHERE id=$3", post.Title, post.Content, postId)
	return err
}

func (s *Storage) DeletePost(id interface{}) error {
	_, err := s.Exec("DELETE FROM posts WHERE id=$1", id)
	return err
}

func (s *Storage) GetUsers() ([]models.User, error) {
	var users []models.User
	err := s.QueryRow("SELECT * FROM users").Scan(&users)
	return users, err
}

func (s *Storage) GetUser(email string) (models.User, error) {
	var user models.User
	err := s.QueryRow("SELECT * FROM users WHERE email=$1", email).Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.Password)
	return user, err
}

func (s *Storage) UserExists(email string) (bool, error) {
	var exists bool
	err := s.QueryRow("SELECT EXISTS(SELECT * FROM users WHERE email=$1)", email).Scan(&exists)
	return exists, err
}

func (s *Storage) AddUser(user models.User, hash []byte) (int, error) {
	var id int
	err := s.QueryRow("INSERT INTO users (email, name, password, phone) VALUES ($1, $2, $3, $4) RETURNING id", user.Email, user.Name, hash, user.Phone).Scan(&id)
	return id, err
}
