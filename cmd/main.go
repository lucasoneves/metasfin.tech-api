package main

import (
	"log"
	"net/http"
	"time"

	"metasfin.tech/controllers"
	"metasfin.tech/database"
	"metasfin.tech/initializers"
	"metasfin.tech/middlewares"
	"metasfin.tech/models"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"
)

func main() {

	initializers.LoadEnvs()
	database.InitDatabase()

	err := database.DB.AutoMigrate(&models.Goal{}, &models.User{})
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

	// Define uma rota GET para o caminho "/".
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Bem-vindo à API de Metas Financeiras! 👋",
		})
	})

	// --- Rotas CRUD para Metas (Goals) ---
	router.POST("/api/auth/signup", controllers.CreateUser)
	router.POST("/api/auth/login", controllers.Login)
	router.GET("/user/profile", middlewares.CheckAuth, controllers.GetUserProfile)
	router.GET("/api/goals", middlewares.CheckAuth, controllers.GetGoals)
	router.GET("/api/goals/:id", middlewares.CheckAuth, controllers.GetGoalByID)
	router.POST("/api/goals", middlewares.CheckAuth, controllers.CreateGoal)
	router.PUT("/api/goals/:id", middlewares.CheckAuth, controllers.UpdateGoal)
	router.DELETE("/api/goals/:id", middlewares.CheckAuth, controllers.DeleteGoal)

	router.POST("/api/goals/deposit/:id", middlewares.CheckAuth, controllers.AddMoneyToGoal)

	router.GET("/api/goals/info", controllers.GetGoalsInfoDashboard)
	log.Printf("Servidor Gin rodando na porta :8080")
	router.Run(":8080")
}
