package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvs() {
	err := godotenv.Load()
	// Se o arquivo .env não for encontrado, não é um erro fatal,
	// pois as variáveis podem ser definidas no ambiente (ex: Docker).
	if err != nil {
		log.Printf("Aviso: Não foi possível carregar o arquivo .env. Usando variáveis de ambiente. Erro: %v", err)
	}

}
