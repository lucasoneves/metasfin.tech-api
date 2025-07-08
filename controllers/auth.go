// controllers/auth_controller.go
package controllers

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	// Ajuste o caminho do import
	"metasfin.tech/database" // Ajuste o caminho do import para o seu Gorm DB

	"metasfin.tech/config"
	"metasfin.tech/models" // Ajuste o caminho do import para o seu modelo de usuário

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// UserInfo representa as informações básicas do usuário retornadas pelo Google
type GoogleUserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	// Adicione outros campos que você possa precisar, como 'picture'
}

// TokenClaims é a estrutura para o seu token JWT
type TokenClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func init() { // Bloco init para carregar a chave JWT ao iniciar o pacote
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("Variável de ambiente JWT_SECRET não definida. Por favor, defina-a.")
	}
	jwtSecret = []byte(secret)
}

// GoogleLogin redireciona para a página de login do Google (usado mais para teste direto)
func GoogleLogin(c *gin.Context) {
	// Gerar um state aleatório para segurança (prevenção de CSRF)
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao gerar o estado de segurança"})
		return
	}
	state := fmt.Sprintf("%x", b)

	// Armazenar o state em um cookie seguro
	c.SetCookie("oauthstate", state, 3600, "/", "localhost", false, true) // 1 hora de validade

	url := config.GoogleOAuthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleAuthCallback lida com o callback do Google
func GoogleAuthCallback(c *gin.Context) {
	// Adicionar verificação defensiva para o banco de dados
	if database.DB == nil {
		log.Println("ERRO CRÍTICO: A conexão com o banco de dados (database.DB) é nula.")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor: conexão com o banco de dados não está disponível."})
		return
	}

	oauthState, _ := c.Cookie("oauthstate")
	state := c.Query("state")

	// Compara o state do cookie com o state da query para previnir ataques CSRF
	if state != oauthState {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Estado inválido"})
		return
	}

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Código de autorização não fornecido"})
		return
	}

	// Troca o código por um token
	token, err := config.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao trocar código por token", "details": err.Error()})
		return
	}

	// Obtém informações do usuário do Google
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao obter informações do usuário do Google", "details": err.Error()})
		return
	}
	defer resp.Body.Close()

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao ler dados do usuário do Google", "details": err.Error()})
		return
	}

	var googleUser GoogleUserInfo
	if err := json.Unmarshal(userData, &googleUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao analisar dados do usuário do Google", "details": err.Error()})
		return
	}

	// Aqui você deve lidar com o usuário: verificar se existe, criar ou atualizar
	var user models.User // Assumindo que você tem um modelo 'User'
	result := database.DB.Where("email = ?", googleUser.Email).First(&user)

	if result.Error != nil { // Usuário não encontrado, criar novo
		user = models.User{
			Name:  googleUser.Name,
			Email: googleUser.Email,
			// Você pode adicionar um campo 'GoogleID' ao seu modelo User para armazenar googleUser.ID
			// GoogleID: googleUser.ID,
		}
		if err := database.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar usuário no banco de dados", "details": err.Error()})
			return
		}
	} else {
		// Usuário encontrado, talvez atualizar informações (nome, etc.)
		user.Name = googleUser.Name
		// user.GoogleID = googleUser.ID
		database.DB.Save(&user)
	}

	// Gerar um JWT para o usuário
	expirationTime := time.Now().Add(24 * time.Hour) // Token válido por 24 horas
	claims := &TokenClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao gerar token JWT", "details": err.Error()})
		return
	}

	// Retorna o token para o frontend
	c.JSON(http.StatusOK, gin.H{
		"message": "Autenticação Google bem-sucedida",
		"token":   tokenString,
		"user":    gin.H{"id": user.ID, "name": user.Name, "email": user.Email},
	})
}
