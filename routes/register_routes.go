package routes

import (
	"github.com/elwafa/billion-data/internal/handlers"
	"github.com/elwafa/billion-data/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	UserHandler *handlers.UserHandler
	AuthHandler *handlers.AuthHandler
	ItemHandler *handlers.ItemHandler
}

func RegisterRoutes(router *gin.Engine, handler *Handler) {
	// handle uploads
	router.Static("/uploads", "./uploads")

	router.POST("/users", handler.UserHandler.StoreUser)

	// auth routes
	router.POST("/login", handler.AuthHandler.Login)

	// groups customer
	customer := router.Group("/customer")
	customer.Use(middleware.AuthMiddleware())
	customerItems := customer.Group("/items")
	customerItems.GET("/", handler.ItemHandler.GetItemsForCustomer)

	// groups seller
	seller := router.Group("/seller")
	seller.Use(middleware.AuthMiddleware())
	items := seller.Group("/items")
	items.POST("/", handler.ItemHandler.StoreItem)
	items.GET("/", handler.ItemHandler.GetItemsForSeller)
}
