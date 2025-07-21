package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"metasfin.tech/database"
	"metasfin.tech/models"
)

func CreateUser(c *gin.Context) {
	var authInput models.AuthInput

	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificação de usuário e email existentes de forma mais robusta
	var existingUser models.User
	// Checa se o username OU o email já existem em uma única consulta
	result := database.DB.Where("username = ? OR email = ?", authInput.Username, authInput.Email).First(&existingUser)
	if result.Error == nil { // Se não deu erro, encontrou um usuário
		if existingUser.Username == authInput.Username {
			c.JSON(http.StatusConflict, gin.H{"error": "Nome de usuário já existe"})
			return
		}
		if existingUser.Email == authInput.Email {
			c.JSON(http.StatusConflict, gin.H{"error": "Email já está em uso"})
			return
		}
	} else if result.Error != gorm.ErrRecordNotFound {
		// Se o erro não for "registro não encontrado", é um erro de banco de dados
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar usuário"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar hash da senha"})
		return
	}

	user := models.User{
		Username: authInput.Username,
		Email:    authInput.Email,
		Password: string(passwordHash),
	}

	if result := database.DB.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar usuário"})
		return
	}

	// Nunca retorne a senha, mesmo que hasheada.
	user.Password = ""
	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func Login(c *gin.Context) {

	var authInput models.LoginInput

	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// Usar First para ter um erro explícito se o usuário não for encontrado
	if err := database.DB.Where("email = ?", authInput.Email).First(&user).Error; err != nil {
		// Resposta genérica para não revelar se o email existe ou não (segurança)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou senha inválidos"})
		return
	}

	// Compara a senha fornecida com o hash armazenado
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authInput.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou senha inválidos"})
		return
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao gerar token"})
		return
	}

	// Nunca retorne a senha na resposta da API
	user.Password = ""
	c.JSON(200, gin.H{
		"token": token,
		"user":  user,
	})
}

func GetUserProfile(c *gin.Context) {

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var user models.User
	// Busca o usuário pelo ID obtido do token. Usamos First para obter um erro claro se não for encontrado.
	// O middleware já garante que userID é um uint.
	result := database.DB.First(&user, userID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar perfil do usuário"})
		}
		return
	}

	user.Password = "" // Garante que a senha nunca seja exposta
	c.JSON(http.StatusOK, gin.H{"user": user})
}
