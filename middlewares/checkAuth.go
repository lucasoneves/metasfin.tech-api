package middlewares

import (
	"fmt"

	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func CheckAuth(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido"})
		c.Abort()
		return
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1) // Extrai o token removendo "Bearer "

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inválido: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil // Substitua "your-secret-key" pela sua chave secreta
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
		c.Abort()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["id"].(float64)) // Extrai o ID do usuário das claims (atenção ao tipo!)
		c.Set("userID", userID)                // Define o userID no contexto
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
		c.Abort()
	}
}
