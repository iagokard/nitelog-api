package routes

import (
	"net/http"
	_ "nitelog/docs"
	"nitelog/internal/config"
	"nitelog/internal/middleware"

	meetingHandler "nitelog/internal/handlers/meeting"
	userHandler "nitelog/internal/handlers/user"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"

	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(router *gin.Engine, client *firestore.Client) {
	cfg := config.Load()

	router.GET("/apidoc", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/apidoc/index.html")
	})
	router.GET("/apidoc/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(middleware.TimeoutMiddleware())

	{
		meetings := router.Group("/meetings")
		meetings.POST("", meetingHandler.CreateMeeting)
		meetings.GET("/by-date/:date", meetingHandler.GetMeetingByDate)
		meetings.GET("/:id", meetingHandler.GetMeetingByID)
		meetings.POST("/add-attendance", meetingHandler.AddUserAttendance)
		meetings.POST("/finish-attendance", meetingHandler.FinishUserAttendance)

		meetings.Use(
			middleware.CORSMiddleware(),
			middleware.JWTMiddleware(cfg.JWTSecret),
		)

		meetings.DELETE("/delete/:id", meetingHandler.DeleteMeeting)
	}

	{
		users := router.Group("/users")
		users.POST("/register", userHandler.CreateUser)
		users.POST("/login", userHandler.LoginUser)
		users.GET("/:id", userHandler.GetUserByID)
		users.GET("/", userHandler.GetUsers)

		users.Use(
			middleware.CORSMiddleware(),
			middleware.JWTMiddleware(cfg.JWTSecret),
		)

		users.DELETE("/delete/:id", userHandler.DeleteUser)
		users.PUT("/update/:id", userHandler.UpdateUser)
	}
}
