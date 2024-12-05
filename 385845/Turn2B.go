package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func main() {
	// Define the signing key (in a real application, use a more secure method)
	signingKey := []byte("your_secret_signing_key")

	// Create a new JWT claim set
	claims := jwt.MapClaims{
		// Registered Claims
		"iss": "your_domain.com",
		"sub": "user123",
		"exp": time.Now().Add(time.Hour * 1).Unix(), // Token expires in 1 hour
		"nbf": time.Now().Unix(),                    // Token not before current time
		"iat": time.Now().Unix(),                    // Token issued at current time
		"aud": "api.your_domain.com",

		// Public Claim
		"user_id": "456",
		"name":    "John Doe",
	}

	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("Generated Token:", signedToken)

	// Validate the token
	validatedClaims, err := validateJWT(signedToken, signingKey)
	if err != nil {
		fmt.Println("Error validating token:", err)
		return
	}

	// Handle claims correctly
	fmt.Println("Valid Token Claims:")
	fmt.Println("Issuer:", validatedClaims["iss"])
	fmt.Println("Subject:", validatedClaims["sub"])
	fmt.Println("User ID:", validatedClaims["user_id"])
	fmt.Println("Name:", validatedClaims["name"])
}

// Define the function to validate the JWT
func validateJWT(tokenString string, signingKey []byte) (jwt.MapClaims, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Check if the token is expired
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		expiration, expOk := claims["exp"].(float64)
		if !expOk {
			return nil, fmt.Errorf("invalid expiration claim")
		}

		expiredTime := time.Unix(int64(expiration), 0)
		if time.Now().After(expiredTime) {
			return nil, fmt.Errorf("token is expired")
		}
	}

	// Return the claims if the token is valid
	return token.Claims.(jwt.MapClaims), nil
}
