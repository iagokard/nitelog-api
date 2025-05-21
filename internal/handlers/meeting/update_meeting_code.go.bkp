package meeting

import (
	"context"
	"errors"
	"net/http"
	"time"

	"nitelog/internal/models"
	"nitelog/internal/util"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *MeetingController) UpdateMeetingCode(c *gin.Context) {
	dateParam := c.Param("date")
	date, err := time.Parse("2006-01-02", dateParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	var meeting models.Meeting
	ctx := context.Background()
	err = h.collection.FindOne(ctx, bson.M{"date": date}).Decode(&meeting)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No meeting for this date"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newCode := util.GenerateMeetingCode()
	_, err = h.collection.UpdateByID(
		ctx,
		meeting.ID,
		bson.M{"$set": bson.M{
			"meeting_code": newCode,
			"updated_at":   time.Now(),
		}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"new_meeting_code": newCode})
}
