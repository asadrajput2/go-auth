package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/asadrajput2/go-auth/pkg/jwt"
	"github.com/asadrajput2/go-auth/pkg/models"
	"github.com/asadrajput2/go-auth/pkg/postgres"
	"golang.org/x/crypto/bcrypt"
)

func Protected(db *postgres.Storage, f func(db *postgres.Storage, userId interface{}, w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
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

func Signup(db *postgres.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Endpoint hit: signup")
		reqBody, _ := ioutil.ReadAll(r.Body)
		var user models.User
		err := json.Unmarshal(reqBody, &user)
		fmt.Println(err)

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
		userExists, err := db.UserExists(user.Email)

		if err != nil {
			log.Fatal(err)
		}
		if userExists {
			json.NewEncoder(w).Encode(map[string]string{"error": "email already exists"})
			return
		}

		// if user doesn't exist, create user -----

		// hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost) // TODO: change cost

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode("{message: Something went wrong on the server}")
			fmt.Println("error hashing password", err)
			return
		}

		// save
		newUserId, err := db.AddUser(user, hash)

		if err != nil {
			log.Fatal(err)
		}

		// generate token
		token, err := jwt.GenerateToken(uint8(newUserId))
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

func Login(db *postgres.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Endpoint hit: login")
		reqBody, err := ioutil.ReadAll(r.Body)

		if err != nil {
			fmt.Println("invalid body")
			json.NewEncoder(w).Encode(map[string]interface{}{"error": err})
			return
		}

		var user models.User
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
		dbUser, err := db.GetUser(user.Email)

		if err != nil || dbUser == (models.User{}) {
			json.NewEncoder(w).Encode(map[string]string{"error": "email not found"})
			return
		}

		// if user exists, check password
		err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))

		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		tokenString, err := jwt.GenerateToken(dbUser.Id)

		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "something went wrong"})
			return
		}
		// return token if user exists and password is correct
		json.NewEncoder(w).Encode(map[string]string{"message": "success", "token": tokenString})

	}
}

func VerifyToken(db *postgres.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Endpoint hit: verify token")
		tokenString := r.Header.Get("Authorization")
		_, err := jwt.ValidateToken(tokenString)
		if err != nil {
			fmt.Fprintf(w, "invalid token")
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "success"})
	}
}
