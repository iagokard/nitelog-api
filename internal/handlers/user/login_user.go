package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"nitelog/internal/config"
	"nitelog/internal/services/user"
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
func LoginUser(c *gin.Context) {
	var req LoginUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()

	userService := services.NewUserService()
	user, err := userService.GetByEmail(ctx, req.Email)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email"})
		return
	}

	if err := util.CheckPassword(user.PasswordHash, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	cfg := config.Load()
	token, err := util.GenerateJWT(user.ID, cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	res := LoginUserResponse{
		Token: token,
	}

	c.JSON(http.StatusOK, res)
}
