package main

import (
	"log"
	"net/http"
	"time"

	"metasfin.tech/controllers" // Import your controllers package
	"metasfin.tech/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	// Adicione esta linha
	"github.com/gin-contrib/cors" // Adicione esta linha
)

// No need for a global DB here in main.go if you're assigning directly to controllers.DB
// var DB *gorm.DB // You can remove this line!

func main() {
	// Call initDatabase to establish the connection
	// and assign it to controllers.DB
	initDatabase()

	// Migra o esquema do banco de dados (cria a tabela 'goals' se não existir).
	// O AutoMigrate adicionará colunas se a struct Goal mudar (com cuidado).
	err := controllers.DB.AutoMigrate(&models.Goal{}) // Use controllers.DB here!
	if err != nil {
		log.Fatalf("Falha ao migrar o banco de dados: %v", err)
	}
	log.Println("Migração do banco de dados concluída com sucesso.")

	// Inicializa um novo "engine" (motor) do Gin.
	// O "Default" inclui middlewares como logger e recovery.
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

	// GET /goals - Lista todas as metas
	router.GET("/api/goals", controllers.GetGoals)
	router.GET("/api/goals/:id", controllers.GetGoalByID)
	router.POST("/api/goals", controllers.CreateGoal)
	router.PUT("/api/goals/:id", controllers.UpdateGoal)
	router.DELETE("/api/goals/:id", controllers.DeleteGoal)

	log.Printf("Servidor Gin rodando na porta :8080")
	router.Run(":8080")
}

// initDatabase inicializa a conexão com o PostgreSQL usando GORM.
func initDatabase() {
	// dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Sao_Paulo",
	//  os.Getenv("DB_HOST"), // Ex: "localhost" quando rodando via Docker na mesma máquina
	//  os.Getenv("DB_USER"),
	//  os.Getenv("DB_PASSWORD"),
	//  os.Getenv("DB_NAME"),
	//  os.Getenv("DB_PORT"),
	// )

	// Para facilitar o desenvolvimento local sem precisar configurar env vars no terminal,
	// você pode usar valores fixos temporariamente (descomente a linha abaixo):
	dsn := "host=localhost user=user password=admin123 dbname=metasfin.tech port=5432 sslmode=disable TimeZone=America/Sao_Paulo"

	var err error
	// Assign the opened DB connection directly to controllers.DB
	controllers.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}) // <--- THIS IS THE KEY CHANGE
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}

	log.Println("Conexão com o banco de dados estabelecida!")
}
