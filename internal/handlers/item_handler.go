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
	color := c.PostForm("color")
	category := c.PostForm("category")
	// convert price to float64
	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// get user id from middleware
	// upload item picture
	userId := c.MustGet("userId").(int)
	item, err := entities.NewItem(name, filePath, description, status, receive, color, category, priceFloat, userId)
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

func (h *ItemHandler) GetItemsForSeller(c *gin.Context) {
	// get itesm from service for seller
	userId := c.MustGet("userId").(int)
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}
	items, total, err := h.service.GetItemsForSeller(c, limit, page, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// add app domain to the picture path
	for i := range items {
		items[i].Picture = h.AppDomain + "/" + items[i].Picture
	}
	// add pagination to the response

	// calculate total pages
	// if total items is not divisible by limit, add 1 to total pages
	totalPages := total % limit
	if totalPages > 0 {
		totalPages = total/limit + 1
	} else {
		totalPages = total / limit
	}
	// claculate last page
	lastPage := page
	if page > 1 {
		lastPage = totalPages
	}
	pagination := map[string]int{
		"current":     page,
		"limit":       limit,
		"total_items": total,
		"total_pages": totalPages,
		"lastPage":    lastPage,
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "pagination": pagination})
}

func (h *ItemHandler) GetItemsForCustomer(c *gin.Context) {
	// get items from service for customer
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}
	sort := c.Query("sort")
	if sort != "ASC" && sort != "DESC" {
		sort = "ASC"
	}
	name := c.Query("name")
	colour := c.Query("colour")
	category := c.Query("category")
	price, err := strconv.ParseFloat(c.Query("price"), 64)
	items, total, err := h.service.GetItemsForCustomer(c, limit, page, sort, name, colour, category, price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// add app domain to the picture path
	for i := range items {
		items[i].Picture = h.AppDomain + "/" + items[i].Picture
	}
	// add pagination to the response

	// calculate total pages
	// if total items is not divisible by limit, add 1 to total pages
	totalPages := total % limit
	if totalPages > 0 {
		totalPages = total/limit + 1
	} else {
		totalPages = total / limit
	}
	// claculate last page
	lastPage := page
	if page > 1 {
		lastPage = totalPages
	}
	pagination := map[string]int{
		"current":     page,
		"limit":       limit,
		"total_items": total,
		"total_pages": totalPages,
		"lastPage":    lastPage,
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "pagination": pagination})
}

func (h *ItemHandler) GetItemForCustomer(c *gin.Context) {
	// get item from service for customer
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.service.GetItemForCustomer(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// add app domain to the picture path
	item.Picture = h.AppDomain + "/" + item.Picture
	c.JSON(http.StatusOK, gin.H{"item": item})

}

func (h *ItemHandler) DeleteItem(c *gin.Context) {
	// delete item from service
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId := c.MustGet("userId").(int)
	err = h.service.DeleteItem(c, id, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})

}
