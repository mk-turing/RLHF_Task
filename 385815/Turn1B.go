package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TenantInfo holds tenant-specific data.
type TenantInfo struct {
	ID       string `json:"id"`
	Secret   string `json:"secret"`
	Name     string `json:"name"`
	Domain   string `json:"domain"`
}

// CustomClaims defines custom JWT claims.
type CustomClaims struct {
	TenantID string `json:"tenant_id"`
	UserID   string `json:"user_id"`
	jwt.StandardClaims
}

var tenants = []TenantInfo{
	{ID: "1", Secret: "super-secret-key-for-tenant-1", Name: "Tenant One", Domain: "tenant1.example.com"},
	{ID: "2", Secret: "another-secret-key-for-tenant-2", Name: "Tenant Two", Domain: "tenant2.example.com"},
}

func findTenantByID(tenantID string) *TenantInfo {
	for _, tenant := range tenants {
		if tenant.ID == tenantID {
			return &tenant
		}
	}
	return nil
}

func generateJWT(tenantID string, userID string) (string, error) {
	tenant := findTenantByID(tenantID)
	if tenant == nil {
		return "", fmt.Errorf("tenant not found: %s", tenantID)
	}

	claims := &CustomClaims{
		TenantID: tenantID,
		UserID:   userID,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tenant.Secret))
}

func verifyJWT(tokenString string) (*CustomClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			return nil, fmt.Errorf("invalid claims")
		}
		tenant := findTenantByID(claims.TenantID)
		if tenant == nil {
			return nil, fmt.Errorf("tenant not found: %s", claims.TenantID)
		}
		return []byte(tenant.Secret), nil
	})

	if claims, ok := token.Claims.(*CustomClaims); !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Token not provided", http.StatusUnauthorized)
		return
	}

	parts := strings.SplitN(tokenString, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "Invalid token format", http.StatusUnauthorized)
		return
	}

	claims, err := verifyJWT(parts[1])
	if err != nil {