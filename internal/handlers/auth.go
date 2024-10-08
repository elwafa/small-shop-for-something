package handlers

import (
	"github.com/elwafa/billion-data/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err, token := h.service.Login(ctx, req.Email, req.Password)
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
		ID:       user.Id,
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		Type:     user.Type,
		IsActive: user.IsActive,
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Login successfully", "token": token, "user": userResponse})
}

func (h *AuthHandler) RenderWebLogin(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", nil)
}

func (h *AuthHandler) WebLogin(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	// bind the request from html form
	user, err := h.service.AdminLogin(ctx, email, password)
	if err != nil {
		// redirect to login page with error message
		ctx.HTML(http.StatusUnauthorized, "login.html", gin.H{"Error": err.Error()})
		return
	}
	// redirect to dashboard
	ctx.HTML(http.StatusFound, "/dashboard", gin.H{"user": user})
}
