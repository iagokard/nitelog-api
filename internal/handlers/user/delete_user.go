package user

import (
	"context"
	"errors"
	"net/http"
	"nitelog/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// DeleteMeeting godoc
// @Summary      Deleta um usuário
// @Description  Deleta um usuário do banco de dados
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user_id   path   string true "Id do usuário (BSON primitive.ObjectID)"
// @Success      200         {object}  util.MessageResponse
// @Failure      400         {object}  util.ErrorResponse
// @Failure      404         {object}  util.ErrorResponse
// @Failure      500         {object}  util.ErrorResponse
// @Router       /users/:id [delete]
func (h *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	ctx := context.Background()

	userService := services.NewUserService()
	err = userService.SoftDelete(ctx, objID)

	if errors.Is(err, mongo.ErrNoDocuments) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meeting not found"})
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
