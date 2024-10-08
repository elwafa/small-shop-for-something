package handlers

import (
	"github.com/elwafa/billion-data/internal/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type CardHandler struct {
	service *services.CardService
}

func NewCardHandler(service *services.CardService) *CardHandler {
	return &CardHandler{service: service}
}

func (h *CardHandler) AddToCard(ctx *gin.Context) {
	// get item id from url and convert it to int
	itemID, err := strconv.Atoi(ctx.Param("item-id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// get user id from context
	userID := ctx.GetInt("userId")
	err = h.service.AddToCard(ctx, itemID, userID)
	log.Println(err, "err")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Item added successfully"})
}
