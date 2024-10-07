package handlers

import (
	"errors"
	"github.com/elwafa/billion-data/internal/entities"
	"github.com/elwafa/billion-data/internal/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type UserHandler struct {
	userService *services.UserService
	authService *services.AuthService
}

func NewUserHandler(service *services.UserService, authService *services.AuthService) *UserHandler {
	return &UserHandler{userService: service, authService: authService}
}

func (h *UserHandler) StoreUser(ctx *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Phone    string `json:"phone"`
		Type     string `json:"type"`
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// validate phone number is valid
	user, err := entities.NewUser(req.Name, req.Email, req.Password, req.Phone, req.Type, true)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.userService.StoreUser(ctx, user)
	if errors.Is(err, services.ErrUserAlreadyExist) {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	loginUser, err, token := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userResponse := struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Type     string `json:"type"`
		IsActive bool   `json:"is_active"`
	}{
		ID:       loginUser.Id,
		Name:     loginUser.Name,
		Email:    loginUser.Email,
		Phone:    loginUser.Phone,
		Type:     loginUser.Type,
		IsActive: loginUser.IsActive,
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User created successfully", "token": token, "user": userResponse})
}
