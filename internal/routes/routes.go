package routes

import (
	_ "nitelog/docs"
	"nitelog/internal/config"
	"nitelog/internal/handlers/meeting"
	"nitelog/internal/handlers/user"
	"nitelog/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(router *gin.Engine, db *mongo.Database) {
	cfg := config.Load()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	meetingHandler := meeting.NewMeetingController(db.Collection("meetings"))
	{
		meetings := router.Group("/meetings")
		meetings.POST("", meetingHandler.CreateMeeting)
		meetings.GET("/by-date/:date", meetingHandler.GetMeetingByDate)
		meetings.POST("/add-attendance", meetingHandler.AddUserAttendance)
		meetings.POST("/finish-attendance", meetingHandler.FinishUserAttendance)
		meetings.GET("/:id", meetingHandler.GetMeetingByID)
		meetings.PUT("/update/:date", meetingHandler.UpdateMeetingCode)
		// meetings.PUT("/:id", handler.UpdateMeeting)
		meetings.DELETE("/:id", meetingHandler.DeleteMeeting)

	}

	userHandler := user.NewUserController(db.Collection("users"))
	{
		users := router.Group("/users")
		users.POST("/register", userHandler.CreateUser)
		users.POST("/login", userHandler.LoginUser)

		users.Use(
			middleware.CORSMiddleware(),
			middleware.JWTMiddleware(cfg.JWTSecret),
		)

		users.PUT("/update/:id", userHandler.UpdateUser)
	}
}
