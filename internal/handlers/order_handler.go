package handlers

import (
	"github.com/elwafa/billion-data/internal/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type OrderHandler struct {
	service     *services.OrderService
	cardService *services.CardService
}

func NewOrderHandler(service *services.OrderService, cardService *services.CardService) *OrderHandler {
	return &OrderHandler{
		service:     service,
		cardService: cardService,
	}
}

func (h *OrderHandler) StoreOrder(ctx *gin.Context) {
	// get user id from context
	userID := ctx.GetInt("userId")
	log.Println("userID", userID)
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
