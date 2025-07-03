package user

import (
	"context"
	"net/http"

	"nitelog/internal/services/user"

	"github.com/gin-gonic/gin"
)

// GetMeetingByID godoc
// @Summary      Retorna todos os usuários
// @Description  Retorna lista com todos os usuários
// @Tags         user_admin
// @Produce      json
// @Success      200         {object}  []models.User
// @Failure      400         {object}  util.ErrorResponse
// @Failure      500         {object}  util.ErrorResponse
// @securityDefinitions.apikey  BearerAuth
// @Router       /users [get]
func GetUsers(c *gin.Context) {
	ctx := context.Background()

	userService := services.NewUserService()
	users, err := userService.GetAllUsers(ctx)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
