package meeting

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FinishUserAttendanceRequest struct {
	UserID string `json:"user_id" binding:"required" example:"68253a5154c3608b34c81d79"`
	Date   string `json:"date" binding:"required" example:"2025-10-26"`
}

// FinishUserAttendance godoc
// @Summary      Finaliza presença em reunião
// @Description  Finaliza a presença usuário do usuário
// @Tags         attendance
// @Accept       json
// @Produce      json
// @Param        attendance     body   FinishUserAttendanceRequest true "Dados da presença"
// @Success      200         {object}  util.MessageResponse
// @Failure      400         {object}  util.ErrorResponse
// @Failure      403         {object}  util.ErrorResponse
// @Failure      404         {object}  util.ErrorResponse
// @Failure      500         {object}  util.ErrorResponse
// @Router       /meetings/finish-attendance [post]
func (h *MeetingController) FinishUserAttendance(c *gin.Context) {
	var req struct {
		UserID string `json:"userId" binding:"required"`
		Date   string `json:"date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Step 1: Check if the attendance has already been finalized
	filter := bson.M{
		"date": date,
		"attendance": bson.M{
			"$elemMatch": bson.M{
				"user_id":  userID,
				"end_time": bson.M{"$ne": nil},
			},
		},
	}

	var existingRecord bson.M
	err = h.collection.FindOne(ctx, filter).Decode(&existingRecord)
	if err == nil {
		// Attendance already finalized
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Attendance already finalized for this user on the specified date",
		})
		return
	} else if err != mongo.ErrNoDocuments {
		// An error occurred during the query
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error checking attendance status",
		})
		return
	}

	// Step 2: Finalize attendance by setting end_time
	update := bson.M{
		"$set": bson.M{
			"attendance.$[elem].end_time": time.Now(),
		},
	}

	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []any{
			bson.M{"elem.user_id": userID, "elem.end_time": bson.M{"$eq": nil}},
		},
	})

	result, err := h.collection.UpdateOne(
		ctx,
		bson.M{"date": date},
		update,
		arrayFilters,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating attendance"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No matching attendance record found or attendance already finalized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendance finalized successfully"})
}
