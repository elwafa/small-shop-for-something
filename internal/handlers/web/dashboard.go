package web

import (
	"github.com/elwafa/billion-data/internal/services"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	userService *services.UserService
}

func NewDashboardHandler(userService *services.UserService) *DashboardHandler {
	return &DashboardHandler{userService: userService}
}

func (h *DashboardHandler) RenderDashboard(ctx *gin.Context) {
	// get user id from context
	userID := ctx.GetInt("userId")
	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		ctx.HTML(500, "login.html", gin.H{"Error": err.Error()})
		return
	}

	ctx.HTML(200, "dashboard.html", gin.H{"user": user})
}
