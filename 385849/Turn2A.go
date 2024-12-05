package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"net/http"
)

type CustomClaims struct {
	jwt.StandardClaims
	Admin   bool `json:"admin"`
	Role    string `json:"role"`
	Username string `json:"username"`
}

// YourHandler is the example endpoint handler
func YourHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Welcome to your protected endpoint!")
}

func JwtMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	secretKey := []byte("yourSecretKey") // Replace this with a secure key

	return func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header
		// In a real application, you might retrieve it from a cookie or different header, based on your needs.
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Token not found", http.StatusUnauthorized)
			return
		}

		// Split out the Bearer prefix
		if len(tokenString) < 7 || tokenString[:7] != "Bearer " {
			http.Error(w, "Token has invalid format", http.StatusUnauthorized)
			return
		}

		tokenString = tokenString[7:]

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHS256); !ok {
				return nil, fmt.Errorf("Unsupported signing method")
			}
			return secretKey, nil
		})

		if err != nil {
			if vert, ok := err.(*jwt.ValidationError); ok {
				if vert.Errors&jwt.ValidationErrorExpired == jwt.ValidationErrorExpired {
					http.Error(w, "Token expired", http.StatusForbidden)
					return
				}
				if vert.Errors&jwt.ValidationErrorSignatureInvalid == jwt.ValidationErrorSignatureInvalid {
					http.Error(w, "Invalid signature", http.StatusForbidden)
					return
				}
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if the token contains custom claims
		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			http.Error(w, "Token contains invalid claims", http.StatusForbidden)
			return
		}

		// You can validate additional custom claims here.
		// For example:
		//if claims.Role != "admin" {
		//    http.Error(w, "Unauthorized", http.StatusForbidden)
		//    return
		//}

		// Extract the username from the token and set it as a context value for downstream handlers
		ctx := r.Context()
		ctx = context.WithValue(ctx, "username", claims.Username)
		r = r.WithContext(ctx)

		handler(w, r)
	}
}

func main() {
	r := mux.NewRouter()
	r.Use(JwtMiddleware)
	r.HandleFunc("/api/example", YourHandler)