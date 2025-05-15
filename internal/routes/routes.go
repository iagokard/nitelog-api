package routes

import (
	"nitelog/internal/config"
	"nitelog/internal/handlers/meeting"
	"nitelog/internal/handlers/user"
	"nitelog/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(router *gin.Engine, db *mongo.Database) {
	cfg := config.Load()

	meetingHandler := meeting.NewMeetingController(db.Collection("meetings"))
	{
		meetings := router.Group("/meetings")
		meetings.POST("", meetingHandler.CreateMeeting)
		meetings.GET("/by-date/:date", meetingHandler.GetMeetingByDate)
		meetings.GET("/:id", meetingHandler.GetMeetingByID)
		meetings.PUT("/update/:date", meetingHandler.UpdateMeetingCode)
		// meetings.PUT("/:id", handler.UpdateMeeting)
		meetings.DELETE("/:id", meetingHandler.DeleteMeeting)

	}

	userHandler := user.NewUserController(db.Collection("users"))
	{
		users := router.Group("/user")
		users.POST("/add-attendance", userHandler.AddUserAttendance)
		users.POST("/finish-attendance", userHandler.FinishUserAttendance)
		users.POST("/register", userHandler.CreateUser)
		users.POST("/login", userHandler.LoginUser)

		users.Use(
			middleware.CORSMiddleware(),
			middleware.JWTMiddleware(cfg.JWTSecret),
		)

		users.PUT("/update/:id", userHandler.UpdateUser)
	}
}
