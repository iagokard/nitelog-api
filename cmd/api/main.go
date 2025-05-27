package main

import (
	"context"
	"log"
	"os"
	"time"

	"net/http"
	"os/signal"

	"nitelog/internal/config"
	"nitelog/internal/routes"
	"nitelog/internal/services"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// @title           NITELog API
// @version         1.0
// @description     API para gestão do Nite
// @contact.name    Suporte NiteLog
// @contact.email   sample@email.com
// @license.name    MIT
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := firestore.NewClient(ctx, cfg.ProjectID)
	if err != nil {
		log.Fatal("Failed to create Firestore client: ", err)
	}
	defer client.Close()

	_, err = client.Collection("test").Doc("test").Get(ctx)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			log.Fatal("Firestore connection check failed: ", err)
		}

		log.Println("Firestore connection verified")
	}

	services.SetFirestoreClient(client)

	router := gin.Default()
	routes.RegisterRoutes(router, client)

	// Graceful shutdown
	srv := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
