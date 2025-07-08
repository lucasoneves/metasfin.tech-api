package main

import (
	"log"
	"net/http"
	"time"

	"metasfin.tech/controllers"
	"metasfin.tech/database"
	"metasfin.tech/models"

	"metasfin.tech/config"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"
)

func main() {
	database.InitDatabase()
	config.InitGoogleOAuth()

	err := database.DB.AutoMigrate(&models.Goal{})
	if err != nil {
		log.Fatalf("Falha ao migrar o banco de dados: %v", err)
	}
	log.Println("Migração do banco de dados concluída com sucesso.")

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	SetupAuthRoutes(router)

	// Define uma rota GET para o caminho "/".
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Bem-vindo à API de Metas Financeiras! 👋",
		})
	})

	// --- Rotas CRUD para Metas (Goals) ---

	router.GET("/api/goals", controllers.GetGoals)
	router.GET("/api/goals/:id", controllers.GetGoalByID)
	router.POST("/api/goals", controllers.CreateGoal)
	router.PUT("/api/goals/:id", controllers.UpdateGoal)
	router.DELETE("/api/goals/:id", controllers.DeleteGoal)

	router.POST("/api/goals/deposit/:id", controllers.AddMoneyToGoal)

	router.GET("/api/goals/info", controllers.GetGoalsInfoDashboard)
	log.Printf("Servidor Gin rodando na porta :8080")
	router.Run(":8080")
}

func SetupAuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.GET("/google/login", controllers.GoogleLogin) // Mais para teste/exemplo
		auth.GET("/google/callback", controllers.GoogleAuthCallback)
	}
}
