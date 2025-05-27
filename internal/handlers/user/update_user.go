package user

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"nitelog/internal/models"
	"nitelog/internal/services/user"
)

type UpdateUserRequest struct {
	Username string `json:"username" example:"username01"`
	Email    string `json:"email" example:"sample@email.com"`
	Password string `jason:"password" example:"safePassword123#"`
}

// UpdateUser godoc
// @Summary      Atualiza um usuário
// @Description  Atualiza um usuário no sistema
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user_id  path      string  true  "Id do usuário"
// @Param        user_data  body      UpdateUserRequest  true  "Novos Dados do Usuário"
// @Success      200   {object}  util.MessageResponse
// @Failure      400   {object}  util.ErrorResponse
// @Failure      401   {object}  util.ErrorResponse
// @Failure      403   {object}  util.ErrorResponse
// @Failure      500   {object}  util.ErrorResponse
// @Security BearerAuth
// @Router       /users/update/:id [put]
func UpdateUser(c *gin.Context) {
	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idParam := c.Param("id")
	if idParam != uid.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot update other user"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: req.Password,
	}

	ctx := context.Background()
	userService := services.NewUserService()
	err := userService.Update(ctx, idParam, user)

	if errors.Is(err, services.ErrNoChangesDetected) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error updating user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, "user updated successfully")
}
