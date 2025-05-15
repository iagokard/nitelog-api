package user

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"nitelog/internal/models"
	"nitelog/internal/util"
)

type CreateUserRequest struct {
	Username string `json:"username" example:"username01"`
	Email    string `json:"email" example:"sample@email.com"`
	Password string `json:"password" binding:"required" example:"safePassword123#"`
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
func (h *UserController) CreateUser(c *gin.Context) {
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	checkCtx, checkCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer checkCancel()

	var existing models.User
	err := h.collection.FindOne(checkCtx, bson.M{
		"$or": []bson.M{
			{"email": req.Email},
			{"username": req.Username},
		},
	}).Decode(&existing)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":   "User already exists",
			"details": existing,
		})
		return
	}

	hash, err := util.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error hashing password",
		})
		return
	}

	newUser := models.User{
		ID:           primitive.NewObjectID(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hash),
		Roles:        []string{},
	}

	insertCtx, insertCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer insertCancel()

	_, err = h.collection.InsertOne(insertCtx, newUser)
	if err != nil {
		log.Printf("Insert error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create meeting",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}
