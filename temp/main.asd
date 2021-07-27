package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"database/sql"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
)

type Storage struct {
	*sql.DB
}

type Article struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"desc"`
	Content     string `json:"content"`
	AuthorId    string `json:"author_id"`
}

type User struct {
	Id       uint8  `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

var Articles []Article

func Connect() (*Storage, error) {

	sqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5432, "asd", "whythis", "gotest")
	db, err := sql.Open("postgres", sqlInfo)

	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()
	return &Storage{db}, err
}

// // create user table
// func (s *Storage) initializeTables() {
// 	// fmt.Println("database created: users")
// 	// s.Exec("DROP TABLE IF EXISTS users")
// 	// s.Exec("DROP TABLE IF EXISTS posts")
// 	stmt_user, err_user := s.Prepare(
// 		`
// 		CREATE TABLE users(
// 		id SERIAL NOT NULL UNIQUE,
// 		name varchar(50),
// 		email varchar(50),
// 		phone varchar(12));
// 	`)

// 	stmt_post, err_post := s.Prepare(
// 		`CREATE TABLE posts(
// 		id SERIAL NOT NULL,
// 		title varchar(250),
// 		description varchar(250),
// 		content varchar(250),
// 		CONSTRAINT fk_user
// 		FOREIGN KEY (id)
// 		REFERENCES users(id));
// 	`)

// 	if err_user != nil || err_post != nil {
// 		log.Fatal(err_user, err_post)
// 	}

// 	_, err_user = stmt_user.Exec()
// 	_, err_post = stmt_post.Exec()

// 	if err_user != nil || err_post != nil {
// 		log.Fatal(err_user, err_post)
// 	} else {
// 		fmt.Println("tables created successfully")
// 	}
// }

func mainHandler(db *Storage) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Endpoint hit: homepage")
		fmt.Fprintf(w, "working!")
	}
}

func generateToken(userId uint8) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
		"iat":     time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret")) // TODO: change secret

	if err != nil {
		return "", err
	}
	return tokenString, nil

}

func validateToken(tokenString string) (interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil // TODO: change secret
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("claims: ", claims)
		return claims["user_id"], nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

func protected(db *Storage, f func(db *Storage, userId interface{}, w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		userId, err := validateToken(tokenString)
		if err != nil {
			fmt.Fprintf(w, "invalid token")
			return
		}

		f(db, userId, w, r)

	}
}

func verifyToken(db *Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Endpoint hit: verify token")
		tokenString := r.Header.Get("Authorization")
		_, err := validateToken(tokenString)
		if err != nil {
			fmt.Fprintf(w, "invalid token")
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "success"})
	}
}

func login(db *Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Endpoint hit: login")
		reqBody, err := ioutil.ReadAll(r.Body)

		if err != nil {
			fmt.Println("invalid body")
			json.NewEncoder(w).Encode(map[string]interface{}{"error": err})
			return
		}
		var user User
		err = json.Unmarshal(reqBody, &user)
		if err != nil {
			fmt.Println("invalid json")
			json.NewEncoder(w).Encode(map[string]interface{}{"error": err})
			return
		}

		// validate fields
		if user.Email == "" || user.Password == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "email or password invalid",
			})
			fmt.Println("invalid email or password")
			return
		}

		// check if user exists
		rows, err := db.Query("SELECT id, email, password FROM users WHERE email = $1", user.Email)
		if err != nil {
			fmt.Println("error making sb query")
			json.NewEncoder(w).Encode(map[string]string{"error": "something went wrong"})
			return
		}
		defer rows.Close()

		// if user exists, check password
		var tempUser User
		for rows.Next() {
			err = rows.Scan(&tempUser.Id, &tempUser.Email, &tempUser.Password)

			if err != nil {
				fmt.Println("error putting back data")
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(tempUser.Password), []byte(user.Password))

			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			tokenString, err := generateToken(tempUser.Id)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"error": "something went wrong"})
				return
			}

			// return token if user exists and password is correct
			json.NewEncoder(w).Encode(map[string]string{"message": "success", "token": tokenString})

		}
		err = rows.Err()
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "something went wrong"})
			return
		}
	}
}

func signup(db *Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := ioutil.ReadAll(r.Body)
		var user User
		err := json.Unmarshal(reqBody, &user)

		if err != nil {
			log.Fatal(err)
		}

		// validate fields
		if user.Email == "" || user.Password == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "email or password invalid",
			})
			fmt.Println("invalid email or password")
			return
		}

		// check if user exists
		rows, err := db.Query("SELECT id, email FROM users WHERE email = $1", user.Email)
		if err != nil {
			fmt.Println("error making sb query")
			json.NewEncoder(w).Encode(map[string]string{"error": "something went wrong"})
			return
		}
		defer rows.Close()
		var tempUser User
		for rows.Next() {
			err = rows.Scan(&tempUser.Id, &tempUser.Email)
			if err != nil {
				fmt.Println("error putting back data")
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			if tempUser.Email == user.Email {
				json.NewEncoder(w).Encode(map[string]string{"error": "email already exists"})
				return
			}
		}
		err = rows.Err()
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "something went wrong"})
			return
		}

		// if user doesn't exist, create user

		// hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost) // TODO: change cost

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode("{message: Something went wrong on the server}")
			fmt.Println("error hashing password", err)
			return
		}

		// save
		stmt := `INSERT INTO users(
			email,
			password,
			name,
			phone
		) VALUES ($1, $2, $3, $4)
		RETURNING id, email, name, phone;`

		rows, err = db.Query(stmt, user.Email, hash, user.Name, user.Phone)

		if err != nil {
			log.Fatal(err)
		}

		var created_user User
		for rows.Next() {
			err = rows.Scan(&created_user.Id, &created_user.Email, &created_user.Name, &created_user.Phone)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}

		// generate token
		token, err := generateToken(created_user.Id)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		}
		// return token
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "success",
			"token":   token,
		})
	}
}

func getAllArticles(db *Storage, userId interface{}, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: all articles")

	var articleList []Article

	stmt := `SELECT * FROM posts WHERE author_id=$1;`
	rows, err := db.Query(stmt, userId)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var article Article
		err := rows.Scan(&article.Id, &article.Title, &article.Description, &article.Content, &article.AuthorId)
		if err != nil {
			log.Fatal(err)
		}
		articleList = append(articleList, article)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(articleList)

}

func getSingleArticle(db *Storage, userId interface{}, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: single article")
	vars := mux.Vars(r)
	id := vars["id"]

	// var result Article
	stmt := `SELECT * FROM posts WHERE author_id=$1 AND id=$2;`
	rows, err := db.Query(stmt, userId, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var article Article
	for rows.Next() {
		err := rows.Scan(&article.Id, &article.Title, &article.Description, &article.Content, &article.AuthorId)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(article)

}

func deleteArticle(db *Storage, userId interface{}, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: delete article")
	vars := mux.Vars(r)
	id := vars["id"]

	stmt := `DELETE FROM posts WHERE id=$1 RETURNING id`

	rows, err := db.Query(stmt, id)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	var deleted_id int

	for rows.Next() {
		err = rows.Scan(&deleted_id)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = rows.Err()

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(deleted_id)
}

func createArticle(db *Storage, userId interface{}, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: create article")
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article Article
	json.Unmarshal(reqBody, &article)

	// TODO: get and add author id
	stmt := `
			INSERT INTO posts (title, description, content, author_id)
			VALUES ($1, $2, $3, $4) RETURNING *;
		`

	rows, err := db.Query(stmt, article.Title, article.Description, article.Content, userId)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var return_article Article
	for rows.Next() {
		err = rows.Scan(&return_article.Id, &return_article.Title, &return_article.Description, &return_article.Content, &return_article.AuthorId)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "success",
		"data": return_article})
}

func updateArticle(db *Storage, userId interface{}, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: update article")
	reqBody, _ := ioutil.ReadAll(r.Body)

	var new_article Article
	json.Unmarshal(reqBody, &new_article)
	vars := mux.Vars(r)
	id := vars["id"]

	stmt := `UPDATE posts SET 
			title=$1,
			description=$2,
			content=$3
			WHERE id=$4
			RETURNING *;
		`

	rows, err := db.Query(stmt, new_article.Title, new_article.Description, new_article.Content, id)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var return_article Article
	for rows.Next() {
		err = rows.Scan(&return_article.Id, &return_article.Title, &return_article.Description, &return_article.Content, &return_article.AuthorId)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = rows.Err()

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(return_article)
}

func reqHandler(db *Storage) {

	myRouter := mux.NewRouter()

	myRouter.HandleFunc("/", mainHandler(db))
	myRouter.HandleFunc("/articles", protected(db, getAllArticles))

	myRouter.HandleFunc("/article", protected(db, createArticle)).Methods("POST")
	myRouter.HandleFunc("/article/{id}", protected(db, deleteArticle)).Methods("DELETE")
	myRouter.HandleFunc("/article/{id}", protected(db, updateArticle)).Methods("PUT")
	myRouter.HandleFunc("/article/{id}", protected(db, getSingleArticle))

	myRouter.HandleFunc("/signup", signup(db)).Methods("POST")
	myRouter.HandleFunc("/login", login(db)).Methods("POST")
	myRouter.HandleFunc("/verifyToken", verifyToken(db)).Methods("POST")

	fmt.Println("Started server at 8080")
	handler := cors.AllowAll().Handler(myRouter)
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func main() {

	db, err := Connect()
	if err != nil {
		log.Fatal(err)
	}
	// db.initializeTables()

	reqHandler(db)
}
