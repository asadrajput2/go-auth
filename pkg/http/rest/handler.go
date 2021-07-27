package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/asadrajput2/go-auth/pkg/postgres"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func ReqHandler(db *postgres.Storage) {

	myRouter := mux.NewRouter()

	// myRouter.HandleFunc("/", mainHandler(db))
	myRouter.HandleFunc("/articles", Protected(db, getAllArticles))

	myRouter.HandleFunc("/article", Protected(db, CreateArticle)).Methods("POST")
	myRouter.HandleFunc("/article/{id}", Protected(db, DeleteArticle)).Methods("DELETE")
	myRouter.HandleFunc("/article/{id}", Protected(db, UpdateArticle)).Methods("PUT")
	myRouter.HandleFunc("/article/{id}", Protected(db, GetSingleArticle))

	myRouter.HandleFunc("/signup", Signup(db)).Methods("POST")
	myRouter.HandleFunc("/login", Login(db)).Methods("POST")
	myRouter.HandleFunc("/verifyToken", VerifyToken(db)).Methods("POST")

	fmt.Println("Started server at 8080")
	handler := cors.AllowAll().Handler(myRouter)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
