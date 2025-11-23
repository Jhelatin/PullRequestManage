package handler

import (
	"PullRequestManage/internal/service"
	"log"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	teamService *service.TeamService
	userService *service.UserService
	prService   *service.PRService
	logger      *log.Logger
}

func NewHandler(ts *service.TeamService, us *service.UserService, ps *service.PRService, logger *log.Logger) *Handler {
	return &Handler{
		teamService: ts,
		userService: us,
		prService:   ps,
		logger:      logger,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	h.registerTeamRoutes(router.Group("/team"))
	h.registerUserRoutes(router.Group("/users"))
	h.registerPRRoutes(router.Group("/pullRequest"))
}
