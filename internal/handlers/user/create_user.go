package user

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"nitelog/internal/services/user"
	"nitelog/internal/util"
)

type CreateUserRequest struct {
	Name         string `json:"name" example:"John Testes" binding:"required"`
	Registration string `json:"registration" example:"8854652123" binding:"required"`
	Email        string `json:"email" example:"sample@email.com" binding:"required"`
	Password     string `json:"password" example:"safePassword123#" binding:"required"`
}

// CreateUser godoc
// @Summary      Cria um novo usuário
// @Description  Cadastra um novo usuário no sistema
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user  body      CreateUserRequest  true  "Dados do Usuário"
// @Success      201   {object}  models.User
// @Failure      400   {object}  util.ErrorResponse
// @Failure      409   {object}  util.ErrorResponse
// @Failure      500   {object}  util.ErrorResponse
// @Router       /users [post]
func CreateUser(c *gin.Context) {
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	hash, err := util.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error hashing password",
		})
		return
	}

	ctx := context.Background()
	userService := services.NewUserService()
	newUser, err := userService.Create(ctx, req.Registration, req.Email, req.Name, hash)

	if errors.Is(err, services.ErrEmailTaken) || errors.Is(err, services.ErrRegistrationTaken) {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error creating user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}
