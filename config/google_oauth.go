// config/google_oauth.go
package config

import (
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	GoogleOAuthConfig *oauth2.Config
)

// O bloco init é executado automaticamente quando o pacote 'config' é importado pela primeira vez.
// Isso garante que a configuração do OAuth esteja pronta antes de qualquer rota tentar usá-la.
func init() {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Fatal("ERRO CRÍTICO: GOOGLE_CLIENT_ID ou GOOGLE_CLIENT_SECRET não estão definidos. Verifique seu arquivo .env ou as configurações do contêiner.")
	}

	GoogleOAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"), // Deve ser a URL do seu frontend para onde o Google redirecionará
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
