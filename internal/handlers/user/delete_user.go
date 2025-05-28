package user

import (
	"context"
	"errors"
	"net/http"
	"nitelog/internal/services/user"
	"slices"

	"github.com/gin-gonic/gin"
)

// DeleteMeeting godoc
// @Summary      Deleta um usuário
// @Description  Deleta um usuário do banco de dados
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user_id   path   string true "Id do usuário"
// @Success      200         {object}  util.MessageResponse
// @Failure      400         {object}  util.ErrorResponse
// @Failure      404         {object}  util.ErrorResponse
// @Failure      500         {object}  util.ErrorResponse
// @Security BearerAuth
// @Router       /users/delete/:id [delete]
func DeleteUser(c *gin.Context) {
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")

	ctx := context.Background()
	userService := services.NewUserService()

	user, err := userService.GetByID(ctx, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !slices.Contains(user.Roles, "admin") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized, user is not admin"})
		return
	}

	err = userService.SoftDelete(ctx, id)

	if errors.Is(err, services.ErrUserNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User marked as deleted successfully",
		"details": "The user has been soft deleted and can be recovered",
	})
}
