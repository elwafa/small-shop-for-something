package handlers

import (
	"github.com/elwafa/billion-data/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OrderHandler struct {
	service     *services.OrderService
	cardService *services.CardService
	appDomain   string
}

func NewOrderHandler(service *services.OrderService, cardService *services.CardService, appDomain string) *OrderHandler {
	return &OrderHandler{
		service:     service,
		cardService: cardService,
		appDomain:   appDomain,
	}
}

func (h *OrderHandler) StoreOrder(ctx *gin.Context) {
	// get user id from context
	userID := ctx.GetInt("userId")
	order, cardId, err := h.service.StoreOrder(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// make the card finished
	err = h.cardService.UpdateCardToBeFinished(ctx, cardId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "order created", "order": order})
}

func (h *OrderHandler) GetOrders(ctx *gin.Context) {
	// get user id from context
	userID := ctx.GetInt("userId")
	orders, err := h.service.GetOrders(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for _, order := range orders {
		orderItems := order.OrderItems
		for _, orderItem := range orderItems {
			orderItem.Item.Picture = h.appDomain + "/" + orderItem.Item.Picture
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"orders": orders})
}
