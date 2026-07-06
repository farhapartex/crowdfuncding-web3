package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func newAuth0KeyFunc(domain string) (jwt.Keyfunc, error) {
	jwksURL := fmt.Sprintf("https://%s/.well-known/jwks.json", domain)

	k, err := keyfunc.NewDefaultCtx(context.Background(), []string{jwksURL})
	if err != nil {
		return nil, err
	}

	return k.Keyfunc, nil
}

func auth0Middleware(keyFunc jwt.Keyfunc, domain, audience string) gin.HandlerFunc {
	issuer := fmt.Sprintf("https://%s/", domain)

	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		tokenString := strings.TrimPrefix(header, "Bearer ")

		claims := jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, keyFunc,
			jwt.WithValidMethods([]string{"RS256"}),
			jwt.WithIssuer(issuer),
			jwt.WithAudience(audience),
		)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("sub", claims.Subject)
		c.Set("auth0Token", tokenString)
		c.Next()
	}
}

func fetchAuth0UserInfo(domain, accessToken string) (email, name string, err error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s/userinfo", domain), nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("userinfo request failed with status %d", resp.StatusCode)
	}

	var payload struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", "", err
	}

	return payload.Email, payload.Name, nil
}
