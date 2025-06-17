package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"metasfin.tech/models"
)

var DB *gorm.DB

func CreateGoal(c *gin.Context) {
	var newGoal models.Goal

	// Tenta fazer o bind do JSON da requisição para a struct newGoal.
	// O Gin cuida de mapear os campos JSON para os campos da sua struct.
	if err := c.ShouldBindJSON(&newGoal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// DB.Create(&newGoal) insere a nova meta no banco de dados.
	// O GORM preenche automaticamente ID, CreatedAt e UpdatedAt para você.
	result := DB.Create(&newGoal)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar meta"})
		return
	}

	// Retorna a meta criada com status 201 Created.
	c.JSON(http.StatusCreated, newGoal)
}

func GetGoals(c *gin.Context) {
	var goals []models.Goal
	result := DB.Find(&goals) // Busca todas as metas
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar metas"})
		return
	}
	c.JSON(http.StatusOK, goals)
}

// getGoalByID retorna uma meta específica pelo ID do banco de dados.
func GetGoalByID(c *gin.Context) {
	// Pega o parâmetro 'id' da URL (ex: /goals/123).
	idParam := c.Param("id")
	// Converte a string do ID para um inteiro sem sinal (uint), que é o tipo do ID da Goal.
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var goal models.Goal
	// DB.First(&goal, id) busca o primeiro registro que corresponde ao ID.
	// O GORM faz automaticamente um "SELECT * FROM goals WHERE id = X LIMIT 1".
	result := DB.First(&goal, id)
	if result.Error != nil {
		// GORM retorna gorm.ErrRecordNotFound se o registro não for encontrado.
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Meta não encontrada"})
		} else {
			// Outros erros (ex: erro de conexão com DB)
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
	// Primeiro, encontra a meta existente pelo ID no banco de dados.
	result := DB.First(&existingGoal, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Meta não encontrada"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar meta existente"})
		}
		return
	}

	var updatedGoalData models.Goal // Struct para receber os dados JSON da atualização.
	if err := c.ShouldBindJSON(&updatedGoalData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingGoal.Title = updatedGoalData.Title
	existingGoal.Description = updatedGoalData.Description
	existingGoal.Value = updatedGoalData.Value
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
	// DB.Delete(&goal, id) realiza um "soft delete" se a struct tiver gorm.DeletedAt.
	// Ele simplesmente preenche o campo DeletedAt em vez de remover o registro.
	result := DB.Delete(&goal, id)
	// É importante converter 'id' (uint64) para uint, que é o tipo do ID na struct GormModel.
	// result := DB.Delete(&models.Goal{}, uint(id)) // Alternativamente, passe o tipo e o ID diretamente.

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar meta"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meta não encontrada para deletar"})
		return
	}

	// Retorna status 204 No Content para indicar sucesso na exclusão sem corpo de resposta.
	c.Status(http.StatusNoContent)
}

// AddMoneyRequest define a estrutura para o corpo da requisição ao adicionar dinheiro a uma meta.
type AddMoneyRequest struct {
	Amount int `json:"amount" binding:"required,gt=0"` // 'binding:"required,gt=0"' garante que o valor seja positivo.
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
	// Busca a meta pelo ID. É importante converter 'id' (uint64) para uint.
	if result := DB.First(&existingGoal, uint(id)); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Meta não encontrada"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar meta"})
		}
		return
	}

	existingGoal.Balance += req.Amount

	// Verifica se a meta foi concluída
	if existingGoal.Balance >= existingGoal.Value {
		existingGoal.Completed = true
		// Opcional: você pode querer definir existingGoal.Active = false aqui
	}

	if result := DB.Save(&existingGoal); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar saldo da meta"})
		return
	}

	c.JSON(http.StatusOK, existingGoal)
}
