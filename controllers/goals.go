package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"metasfin.tech/database"
	"metasfin.tech/models"
)

type AddMoneyRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

func CreateGoal(c *gin.Context) {
	var newGoal models.Goal

	if err := c.ShouldBindJSON(&newGoal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// MELHORIA DE SEGURANÇA:
	// O UserID NUNCA deve vir do corpo da requisição.
	// Ele deve ser obtido do contexto, onde o middleware de autenticação o coloca.
	userID, exists := c.Get("userID") // Assumindo que o middleware armazena o ID do usuário com a chave "userID"
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}
	newGoal.UserID = userID.(uint) // Converte para o tipo correto (ajuste se necessário)
	result := database.DB.Create(&newGoal)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar meta"})
		return
	}
	c.JSON(http.StatusCreated, newGoal)
}

func GetGoals(c *gin.Context) {
	var goals []models.Goal
	result := database.DB.Find(&goals)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar metas"})
		return
	}
	c.JSON(http.StatusOK, goals)
}

func GetGoalByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var goal models.Goal
	result := database.DB.First(&goal, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Meta não encontrada"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar meta"})
		}
		return
	}
	c.JSON(http.StatusOK, goal)
}

func UpdateGoal(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var existingGoal models.Goal
	result := database.DB.First(&existingGoal, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Meta não encontrada"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar meta existente"})
		}
		return
	}

	var updatedGoalData models.Goal
	if err := c.ShouldBindJSON(&updatedGoalData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingGoal.Title = updatedGoalData.Title
	existingGoal.Description = updatedGoalData.Description
	existingGoal.Balance = updatedGoalData.Balance
	existingGoal.TargetValue = updatedGoalData.TargetValue
	existingGoal.UserID = updatedGoalData.UserID

	result = database.DB.Save(&existingGoal)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar meta"})
		return
	}

	c.JSON(http.StatusOK, existingGoal)
}

func DeleteGoal(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var goal models.Goal
	result := database.DB.Delete(&goal, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar meta"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meta não encontrada para deletar"})
		return
	}
	c.Status(http.StatusNoContent)
}

func AddMoneyToGoal(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req AddMoneyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	var existingGoal models.Goal
	if result := database.DB.First(&existingGoal, uint(id)); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Meta não encontrada"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar meta"})
		}
		return
	}

	existingGoal.Balance += req.Amount

	if existingGoal.Balance >= existingGoal.TargetValue {
		existingGoal.Completed = true
	}

	if result := database.DB.Save(&existingGoal); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar saldo da meta"})
		return
	}

	c.JSON(http.StatusOK, existingGoal)
}

func GetGoalsInfoDashboard(c *gin.Context) {
	var count int64
	var totalBalance float64
	if err := database.DB.Model(&models.Goal{}).Count(&count).Error; err != nil {
		log.Printf("Error counting goals: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao contar metas"})
		return
	}

	// Use .Scan() em vez de .Row().Scan() para melhor tratamento de erros e valores nulos
	if err := database.DB.Model(&models.Goal{}).Select("COALESCE(SUM(balance), 0)").Scan(&totalBalance).Error; err != nil {
		log.Printf("Erro ao calcular o saldo total: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao calcular saldo total"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "Sucesso",
		"total_challenges": count,
		"total_balance":    totalBalance,
	})
}
