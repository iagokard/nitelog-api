package user

import (
	"context"
	"errors"
	"net/http"

	"nitelog/internal/services/user"

	"github.com/gin-gonic/gin"
)

// GetMeetingByID godoc
// @Summary      Procura usuário por id
// @Description  Procura usuário por id especifico
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user_id   path     string true "ID do usuário"
// @Success      200         {object}  models.User
// @Failure      400         {object}  util.ErrorResponse
// @Failure      404         {object}  util.ErrorResponse
// @Failure      500         {object}  util.ErrorResponse
// @securityDefinitions.apikey  BearerAuth
// @Router       /users/:id [get]
func GetUserByID(c *gin.Context) {
	id := c.Param("id")

	ctx := context.Background()

	userService := services.NewUserService()
	user, err := userService.GetByID(ctx, id)

	if errors.Is(err, services.ErrUserNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
