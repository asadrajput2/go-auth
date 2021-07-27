package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/asadrajput2/go-auth/pkg/models"
	"github.com/asadrajput2/go-auth/pkg/postgres"
	"github.com/gorilla/mux"
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

func GetSingleArticle(db *postgres.Storage, userId interface{}, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: single article")
	vars := mux.Vars(r)
	id := vars["id"]

	// var result Article
	article, err := db.GetPost(id)

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(article)

}

func DeleteArticle(db *postgres.Storage, userId interface{}, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: delete article")
	vars := mux.Vars(r)
	id := vars["id"]

	err := db.DeletePost(id)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "success"})
}

func CreateArticle(db *postgres.Storage, userId interface{}, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: create article")
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article models.Article
	json.Unmarshal(reqBody, &article)

	err := db.CreatePost(article, userId)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "success"})
}

func UpdateArticle(db *postgres.Storage, userId interface{}, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: update article")
	reqBody, _ := ioutil.ReadAll(r.Body)

	var new_article models.Article
	json.Unmarshal(reqBody, &new_article)
	vars := mux.Vars(r)
	id := vars["id"]

	err := db.UpdatePost(new_article, id)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "success"})
}
