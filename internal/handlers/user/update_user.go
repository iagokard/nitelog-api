package user

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"nitelog/internal/models"
)

func (h *UserController) UpdateUser(c *gin.Context) {
	// Extrai userID do token
	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Verifica se path param coincide com token
	idParam := c.Param("id")
	if idParam != uid.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot update other user"})
		return
	}

	var req struct {
		Username *string `json:"username" binding:"omitempty,min=3"`
		Email    *string `json:"email" binding:"omitempty,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	oid, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"username": req.Username,
			"email":    req.Email,
		},
	}

	var updated models.User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = h.collection.FindOneAndUpdate(ctx, bson.M{"_id": oid}, update, nil).Decode(&updated)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}
