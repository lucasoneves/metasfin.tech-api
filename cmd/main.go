package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"metasfin.tech/controllers"
	"metasfin.tech/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gin-contrib/cors"
)

func main() {
	initDatabase()

	err := controllers.DB.AutoMigrate(&models.Goal{})
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

func initDatabase() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Sao_Paulo",
		os.Getenv("DATABASE_HOST"), // Será 'db' quando rodando via Docker
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_PORT"),
	)

	var err error
	// Assign the opened DB connection directly to controllers.DB
	controllers.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}) // <--- THIS IS THE KEY CHANGE
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}

	log.Println("Conexão com o banco de dados estabelecida!")
}
