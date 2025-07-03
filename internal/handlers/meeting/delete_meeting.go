package meeting

import (
	"context"
	"errors"
	"net/http"
	meetingServices "nitelog/internal/services/meeting"
	userServices "nitelog/internal/services/user"
	"slices"

	"github.com/gin-gonic/gin"
)

// DeleteMeeting godoc
// @Summary      Deleta uma reunião
// @Description  Deleta uma reunião do banco de dados
// @Tags         meeting_admin
// @Accept       json
// @Produce      json
// @Param        meeting_id   path   string true "Id da reunião"
// @Success      200         {object}  util.MessageResponse
// @Failure      400         {object}  util.ErrorResponse
// @Failure      404         {object}  util.ErrorResponse
// @Failure      500         {object}  util.ErrorResponse
// @securityDefinitions.apikey  BearerAuth
// @Router       /meetings/delete/:id [delete]
func DeleteMeeting(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ctx := context.Background()

	userService := userServices.NewUserService()
	user, err := userService.GetByID(ctx, userID.(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !slices.Contains(user.Roles, "admin") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized, user is not admin"})
		return
	}

	id := c.Param("id")

	meetingService := meetingServices.NewMeetingService()
	err = meetingService.SoftDelete(ctx, id)

	if errors.Is(err, meetingServices.ErrMeetingNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meeting not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Meeting marked as deleted successfully",
		"details": "The meeting has been soft deleted and can be recovered",
	})
}
