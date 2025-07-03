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

	router.Use(middleware.TimeoutMiddleware())
	router.Use(middleware.CORS())

	router.GET("/apidoc", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/apidoc/index.html")
	})
	router.GET("/apidoc/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	{
		meetings := router.Group("/meetings")

		meetings.Use(
			middleware.JWT(cfg.JWTSecret),
		)

		meetings.GET("/by-date/:date", meetingHandler.GetMeetingByDate)
		meetings.GET("/:id", meetingHandler.GetMeetingByID)
		meetings.POST("", meetingHandler.CreateMeeting)
		meetings.POST("/add-attendance", meetingHandler.AddUserAttendance)
		meetings.POST("/finish-attendance", meetingHandler.FinishUserAttendance)

		meetings.Use(middleware.AdminOnly())

		meetings.DELETE("/delete/:id", meetingHandler.DeleteMeeting)
	}

	{
		users := router.Group("/users")
		users.POST("/register", userHandler.CreateUser)
		users.POST("/login", userHandler.LoginUser)

		users.Use(
			middleware.JWT(cfg.JWTSecret),
		)

		users.GET("/:id", userHandler.GetUserByID)
		users.DELETE("/delete/:id", userHandler.DeleteUser)
		users.PUT("/update/:id", userHandler.UpdateUser)

		users.Use(middleware.AdminOnly())

		users.GET("/", userHandler.GetUsers)
	}
}
