package handler

import (
	"PullRequestManage/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) registerTeamRoutes(rg *gin.RouterGroup) {
	rg.POST("/add", h.addTeam)
	rg.GET("/get", h.getTeam)
}

func (h *Handler) addTeam(c *gin.Context) {
	var input domain.Team
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	team, err := h.teamService.CreateTeam(c.Request.Context(), input)
	if err != nil {
		h.errorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"team": team})
}

func (h *Handler) getTeam(c *gin.Context) {
	teamName := c.Query("team_name")
	if teamName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "team_name query parameter is required"})
		return
	}

	team, err := h.teamService.GetTeam(c.Request.Context(), teamName)
	if err != nil {
		h.errorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, team)
}
