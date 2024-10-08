package routes

import (
	"github.com/elwafa/billion-data/internal/handlers"
	"github.com/elwafa/billion-data/internal/handlers/web"
	"github.com/elwafa/billion-data/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	UserHandler      *handlers.UserHandler
	AuthHandler      *handlers.AuthHandler
	ItemHandler      *handlers.ItemHandler
	CardHandler      *handlers.CardHandler
	OrderHandler     *handlers.OrderHandler
	DashboardHandler *web.DashboardHandler
}

func RegisterRoutes(router *gin.Engine, handler *Handler) {
	// handle static files
	router.Static("/assets", "./assets")
	// templates
	router.LoadHTMLGlob("templates/*")
	// Guest Group
	guest := router.Group("/")
	guest.Use(middleware.WebGuest())
	// handle web login
	guest.GET("/", handler.AuthHandler.RenderWebLogin)
	guest.POST("/", handler.AuthHandler.WebLogin)

	// handle Dashboard pages
	dashboard := router.Group("/dashboard")
	dashboard.Use(middleware.WebAuth())
	dashboard.GET("/", handler.DashboardHandler.RenderDashboard)
	//dashboard.GET("/", handler.DashboardHandler.RenderDashboard)
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
	customerItems.GET("/:id", handler.ItemHandler.GetItemForCustomer)
	// customer card
	customerCard := customer.Group("/card")
	customerCard.POST("/:item-id", handler.CardHandler.AddToCard)
	customerCard.GET("/", handler.CardHandler.GetCard)
	// customer order
	customerOrder := customer.Group("/order")
	customerOrder.POST("/", handler.OrderHandler.StoreOrder)
	customerOrder.GET("/", handler.OrderHandler.GetOrders)

	// groups seller
	seller := router.Group("/seller")
	seller.Use(middleware.AuthMiddleware())
	items := seller.Group("/items")
	items.POST("/", handler.ItemHandler.StoreItem)
	items.GET("/", handler.ItemHandler.GetItemsForSeller)
}
