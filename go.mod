module github.com/asadrajput2/go-auth

go 1.16

require (
	github.com/golang-jwt/jwt v3.2.1+incompatible
	github.com/gorilla/mux v1.8.0
	github.com/lib/pq v1.10.2
	github.com/rs/cors v1.8.0
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
)

replace github.com/asadrajput2/go-auth/package => ./pkg
