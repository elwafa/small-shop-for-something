package app

import (
	"github.com/elwafa/billion-data/config"
	"github.com/elwafa/billion-data/internal/db/sql"
	"github.com/elwafa/billion-data/internal/handlers"
	"github.com/elwafa/billion-data/internal/repositories/postgres"
	"github.com/elwafa/billion-data/internal/services"
	"github.com/elwafa/billion-data/routes"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run() {

	log := log.New(os.Stdout, "billion-data ", log.LstdFlags|log.Lshortfile)

	cfg := config.LoadConfig()

	postgresDb, err := sql.InitPostgres(cfg.PostgresDSN)

	if err != nil {
		log.Fatalf("Error initializing postgres: %v", err)
	}
	// register the user repository
	userRepo := postgres.NewPostgresUserRepository(postgresDb)
	itemRepo := postgres.NewPostgresItemRepository(postgresDb)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo)
	itemService := services.NewItemService(itemRepo)
	userHandler := handlers.NewUserHandler(userService, authService)
	authHandler := handlers.NewAuthHandler(authService)
	itemHandler := handlers.NewItemHandler(itemService, cfg.APPDomain)

	h := &routes.Handler{
		UserHandler: userHandler,
		AuthHandler: authHandler,
		ItemHandler: itemHandler,
	}

	router := gin.Default()

	routes.RegisterRoutes(router, h)

	go func() {
		if err := router.Run(":8080"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	gracefulShutdown()
}

func gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gracefully...")

	sql.ClosePostgres()
	log.Println("All connections closed.")
}
