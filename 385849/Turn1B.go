package main

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)


	jwt.StandardClaims
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
	Admin  bool     `json:"admin"`

// generateJWT creates a new JWT token with custom claims
func generateJWT(userID string, roles []string) (string, error) {
	claims := &CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token valid for 24 hours
			Issuer:    "your_issuer",
			Subject:   userID,
		},
		UserID: userID,
		Roles:  roles,
		Admin:  false, // Default: user is not an admin
	}

	secretKey := []byte("your_super_secret_key") // Use a secure key in production
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// parseJWT verifies the token and extracts custom claims
func parseJWT(tokenString string) (*CustomClaims, error) {
	secretKey := []byte("your_super_secret_key")
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}


	claims, err := parseJWT(tokenString)
	if err != nil {
		fmt.Println("Error parsing token:", err)
		return
	}

	fmt.Println("Parsed Claims:")
	fmt.Println("UserID:", claims.UserID)
	fmt.Println("Roles:", claims.Roles)
	fmt.Println("Is Admin:", claims.Admin)

	// Check authorization for a role
	requiredRole := "admin"
	isAuthorized, err := CheckAuthorization(tokenString, requiredRole)
	if err != nil {
		fmt.Println("Authorization check failed:", err)
		return
	}
	if isAuthorized {
		fmt.Printf("User has the '%s' role.\n", requiredRole)
	} else {
		fmt.Printf("User does not have the '%s' role.\n", requiredRole)
	}
}
