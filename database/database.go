package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Sao_Paulo",
		os.Getenv("DATABASE_HOST"), // Será 'db' quando rodando via Docker
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_PORT"),
	)

	var err error
	// Assign the opened DB connection directly to controllers.DB
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}) // <--- THIS IS THE KEY CHANGE
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}

	log.Println("Conexão com o banco de dados estabelecida!")
}
