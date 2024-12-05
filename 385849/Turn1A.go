package main

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

// CustomClaims represent our custom claims for the JWT.
type CustomClaims struct {
	jwt.StandardClaims
	Admin bool   `json:"admin"`
	Role  string `json:"role"`
}

func main() {
	// Generate a secret key for signing the JWT. In a production environment, use a secure key management practice.
	secretKey := []byte("yourSecretKey")

	// Create a new claim.
	claims := CustomClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "your-app-name",
			Subject:   "an example user",
			Audience:  "your-app-audience",
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
		Admin: true,
		Role:  "user",
	}

	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Verify and parse the token
	//parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	//	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	//		return nil, fmt.Errorf("Unsupported signing method")
	//	}
	//	return secretKey, nil
	//})

	parsedToken, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unsupported signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	// Retrieve the custom claims
	customClaims, ok := parsedToken.Claims.(*CustomClaims)
	if ok {
		fmt.Println("User ID:", customClaims.Subject)
		fmt.Println("Is Admin:", customClaims.Admin)
		fmt.Println("Role:", customClaims.Role)
	}
}
