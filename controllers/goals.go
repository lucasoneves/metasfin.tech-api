package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"metasfin.tech/models"
)

var DB *gorm.DB

type AddMoneyRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

func CreateGoal(c *gin.Context) {
	var newGoal models.Goal

	if err := c.ShouldBindJSON(&newGoal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := DB.Create(&newGoal)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar meta"})
		return
	}
	c.JSON(http.StatusCreated, newGoal)
}

func GetGoals(c *gin.Context) {
	var goals []models.Goal
	result := DB.Find(&goals)
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
	result := DB.First(&goal, id)
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
	result := DB.First(&existingGoal, id)
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

	result = DB.Save(&existingGoal)
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
	result := DB.Delete(&goal, id)
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
	if result := DB.First(&existingGoal, uint(id)); result.Error != nil {
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

	if result := DB.Save(&existingGoal); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar saldo da meta"})
		return
	}

	c.JSON(http.StatusOK, existingGoal)
}

func GetGoalsInfoDashboard(c *gin.Context) {
	var count int64
	var totalBalance float64
	result := DB.Model(&models.Goal{}).Count(&count)
	balanceResult := DB.Model(&models.Goal{}).Select("SUM(balance)").Row().Scan(&totalBalance)

	if result.Error != nil {
		log.Fatalf("Error counting goals: %v", result.Error)
	}

	if balanceResult != nil {
		log.Fatalf("Erro ao escanear o resultado da soma: %v", balanceResult)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "Sucesso",
		"total_challenges": count,
		"total_balance":    totalBalance,
	})
}
