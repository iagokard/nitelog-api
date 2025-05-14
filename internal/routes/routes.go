package routes

import (
	"nitelog/internal/config"
	"nitelog/internal/handlers"
	"nitelog/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterMeetingRoutes(router *gin.Engine, db *mongo.Database) {
	cfg := config.Load()

	handler := handlers.NewDataHandler(db, "meetings")
	{
		meetings := router.Group("/meetings")
		meetings.POST("", handler.CreateMeeting)
		meetings.POST("/user", handler.AddUserAttendance)
		meetings.POST("/user/finish", handler.FinishUserAttendance)
		meetings.GET("/by-date/:date", handler.GetMeetingByDate)
		meetings.GET("/:id", handler.GetMeetingByID)
		meetings.PUT("/update/:date", handler.UpdateMeetingCode)
		// meetings.PUT("/:id", handler.UpdateMeeting)
		meetings.DELETE("/:id", handler.DeleteMeeting)

	}

	handler = handlers.NewDataHandler(db, "users")
	{
		user := router.Group("/user")
		user.POST("/register", handler.CreateUser)
		user.POST("/login", handler.LoginUser)

		user.Use(
			middleware.CORSMiddleware(),
			middleware.JWTMiddleware(cfg.JWTSecret),
		)

		user.PUT("/update/:id", handler.UpdateUser)
	}

}
