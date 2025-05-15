package user

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"nitelog/internal/config"
	"nitelog/internal/models"
	"nitelog/internal/util"
)

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginUserResponse struct {
	Token string `json:"token"`
}

// LoginUser godoc
// @Summary      Autentica um usuário
// @Description  Gera token JWT para usuário válido
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        credentials  body      LoginUserRequest  true  "Credenciais"
// @Success      200          {object}  LoginUserResponse
// @Failure      400          {object}  util.ErrorResponse
// @Failure      401          {object}  util.ErrorResponse
// @Failure      500          {object}  util.ErrorResponse
// @Router       /users/login [post]
func (h *UserController) LoginUser(c *gin.Context) {
	var req LoginUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := h.collection.FindOne(ctx, bson.M{
		"email": req.Email,
	}).Decode(&user)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credentials invalid"})
		return
	}

	if err := util.CheckPassword(user.PasswordHash, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credentials invalid"})
		return
	}

	cfg := config.Load()
	token, err := util.GenerateJWT(user.ID.Hex(), cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	res := LoginUserResponse{
		Token: token,
	}

	c.JSON(http.StatusOK, res)
}
