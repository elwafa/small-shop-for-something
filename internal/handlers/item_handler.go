package handlers

import (
	"github.com/elwafa/billion-data/internal/entities"
	"github.com/elwafa/billion-data/internal/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type ItemHandler struct {
	service   *services.ItemService
	AppDomain string
}

func NewItemHandler(service *services.ItemService, appDomain string) *ItemHandler {
	return &ItemHandler{service: service, AppDomain: appDomain}
}

func (h *ItemHandler) StoreItem(c *gin.Context) {
	// upload item picture
	file, err := c.FormFile("picture")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// change file name
	file.Filename = time.Now().Format("2006-01-02 15:04:05") + "_item_" + file.Filename
	// save file to server
	err = c.SaveUploadedFile(file, "uploads/"+file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// get file path to save to database
	filePath := "uploads/" + file.Filename

	name := c.PostForm("name")
	price := c.PostForm("price")
	// print the price and type
	log.Printf("price: %s, type: %T", price, price)
	description := c.PostForm("description")
	status := c.PostForm("status")
	receive := c.PostForm("receive")
	// convert price to float64
	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// get user id from middleware
	// upload item picture
	userId := c.MustGet("userId").(int)

	item, err := entities.NewItem(name, filePath, description, status, receive, priceFloat, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err = h.service.StoreItem(c, item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// add app domain to the picture path
	item.Picture = h.AppDomain + "/" + item.Picture
	c.JSON(http.StatusOK, gin.H{"message": "Item stored successfully", "item": item})

}
