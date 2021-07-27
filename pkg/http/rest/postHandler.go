package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/asadrajput2/go-auth/pkg/models"
	"github.com/asadrajput2/go-auth/pkg/postgres"
)

func getAllArticles(db *postgres.Storage, userId interface{}, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: all articles")

	var articleList []models.Article

	articleList, err := db.GetPosts(10, userId)

	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(articleList)

}

// func getSingleArticle(db *postgres.Storage, userId interface{}, w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Endpoint hit: single article")
// 	vars := mux.Vars(r)
// 	id := vars["id"]

// 	// var result Article
// 	stmt := `SELECT * FROM posts WHERE author_id=$1 AND id=$2;`
// 	rows, err := db.Query(stmt, userId, id)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
// 	var article Article
// 	for rows.Next() {
// 		err := rows.Scan(&article.Id, &article.Title, &article.Description, &article.Content, &article.AuthorId)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// 	err = rows.Err()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	json.NewEncoder(w).Encode(article)

// }

// func deleteArticle(db *Storage, userId interface{}, w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Endpoint hit: delete article")
// 	vars := mux.Vars(r)
// 	id := vars["id"]

// 	stmt := `DELETE FROM posts WHERE id=$1 RETURNING id`

// 	rows, err := db.Query(stmt, id)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer rows.Close()
// 	var deleted_id int

// 	for rows.Next() {
// 		err = rows.Scan(&deleted_id)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}

// 	err = rows.Err()

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	json.NewEncoder(w).Encode(deleted_id)
// }

// func createArticle(db *Storage, userId interface{}, w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Endpoint hit: create article")
// 	reqBody, _ := ioutil.ReadAll(r.Body)

// 	var article Article
// 	json.Unmarshal(reqBody, &article)

// 	// TODO: get and add author id
// 	stmt := `
// 			INSERT INTO posts (title, description, content, author_id)
// 			VALUES ($1, $2, $3, $4) RETURNING *;
// 		`

// 	rows, err := db.Query(stmt, article.Title, article.Description, article.Content, userId)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer rows.Close()

// 	var return_article Article
// 	for rows.Next() {
// 		err = rows.Scan(&return_article.Id, &return_article.Title, &return_article.Description, &return_article.Content, &return_article.AuthorId)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}

// 	err = rows.Err()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]interface{}{"message": "success",
// 		"data": return_article})
// }

// func updateArticle(db *Storage, userId interface{}, w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Endpoint hit: update article")
// 	reqBody, _ := ioutil.ReadAll(r.Body)

// 	var new_article Article
// 	json.Unmarshal(reqBody, &new_article)
// 	vars := mux.Vars(r)
// 	id := vars["id"]

// 	stmt := `UPDATE posts SET
// 			title=$1,
// 			description=$2,
// 			content=$3
// 			WHERE id=$4
// 			RETURNING *;
// 		`

// 	rows, err := db.Query(stmt, new_article.Title, new_article.Description, new_article.Content, id)

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()

// 	var return_article Article
// 	for rows.Next() {
// 		err = rows.Scan(&return_article.Id, &return_article.Title, &return_article.Description, &return_article.Content, &return_article.AuthorId)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}

// 	err = rows.Err()

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	json.NewEncoder(w).Encode(return_article)
// }
