package handlers

import (
	"errors"
	"github.com/elwafa/billion-data/internal/entities"
	"github.com/elwafa/billion-data/internal/repositories"
	"github.com/elwafa/billion-data/internal/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type CardHandler struct {
	service   *services.CardService
	appDomain string
}

func NewCardHandler(service *services.CardService, appDomain string) *CardHandler {
	return &CardHandler{service: service, appDomain: appDomain}
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

func (h *CardHandler) GetCard(ctx *gin.Context) {
	userID := ctx.GetInt("userId")
	card, err := h.service.GetCard(ctx, userID)
	if errors.Is(err, repositories.ErrRecodeNotFound) {
		ctx.JSON(http.StatusOK, gin.H{"card": entities.Card{}})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// handle Item Picture
	items := card.Items
	for _, item := range items {
		item.Item.Picture = h.appDomain + "/" + item.Item.Picture
		log.Println(item.Item.Picture)
	}
	ctx.JSON(http.StatusOK, card)
}
