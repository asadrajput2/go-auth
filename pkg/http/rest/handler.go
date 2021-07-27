package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/asadrajput2/go-auth/pkg/jwt"
	"github.com/asadrajput2/go-auth/pkg/postgres"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func protected(db *postgres.Storage, f func(db *postgres.Storage, userId interface{}, w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		userId, err := jwt.ValidateToken(tokenString)
		if err != nil {
			fmt.Fprintf(w, "invalid token")
			return
		}

		f(db, userId, w, r)

	}
}

func ReqHandler(db *postgres.Storage) {

	myRouter := mux.NewRouter()

	// myRouter.HandleFunc("/", mainHandler(db))
	myRouter.HandleFunc("/articles", protected(db, getAllArticles))

	// myRouter.HandleFunc("/article", protected(db, createArticle)).Methods("POST")
	// myRouter.HandleFunc("/article/{id}", protected(db, deleteArticle)).Methods("DELETE")
	// myRouter.HandleFunc("/article/{id}", protected(db, updateArticle)).Methods("PUT")
	// myRouter.HandleFunc("/article/{id}", protected(db, getSingleArticle))

	// myRouter.HandleFunc("/signup", signup(db)).Methods("POST")
	// myRouter.HandleFunc("/login", login(db)).Methods("POST")
	// myRouter.HandleFunc("/verifyToken", verifyToken(db)).Methods("POST")

	fmt.Println("Started server at 8080")
	handler := cors.AllowAll().Handler(myRouter)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
