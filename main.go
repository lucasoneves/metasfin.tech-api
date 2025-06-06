package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"metasfin.tech/models"
)

func main() {
	fmt.Println("Hello Go!")

	// Inicializa um novo engine do Gin. O "Default" inclui middlewares como logger e recovery

	router := gin.Default()

	// Define uma rota GET para o caminho "api/greeting"
	router.GET("/api/greeting", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Bem-vindo à Metasfin.tech! 👋",
		})
	})

	var myGoal = models.Goal{
		Title:       "Férias de verão 2025",
		Description: "lorem ipsum dolor sit amet",
		Value:       10000,
		UserID:      10,
	}

	var myGoal2 = models.Goal{
		Title:       "Mundial de Clubes 2025",
		Description: "lorem ipsum dolor sit amet",
		Value:       30000,
		UserID:      5,
	}

	// Define uma rota GET para a url "api/goals"
	router.GET("/api/goals", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"data": []models.Goal{myGoal, myGoal2},
		})
	})

	// Inicia o servidor na porta padrão 8080
	router.Run(":8080")
}
