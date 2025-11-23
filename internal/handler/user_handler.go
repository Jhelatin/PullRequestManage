package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) registerUserRoutes(rg *gin.RouterGroup) {
	rg.POST("/setIsActive", h.setIsActive)
	rg.GET("/getReview", h.getReviews)
}

func (h *Handler) setIsActive(c *gin.Context) {
	var input setIsActiveRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	user, err := h.userService.SetUserIsActive(c.Request.Context(), input.UserID, input.IsActive)
	if err != nil {
		h.errorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *Handler) getReviews(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id query parameter is required"})
		return
	}

	prs, err := h.userService.GetReviewsForUser(c.Request.Context(), userID)
	if err != nil {
		h.errorResponse(c, err)
		return
	}

	shortPRs := make([]gin.H, len(prs))
	for i, pr := range prs {
		shortPRs[i] = gin.H{
			"pull_request_id":   pr.ID,
			"pull_request_name": pr.Name,
			"author_id":         pr.AuthorID,
			"status":            pr.Status,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":       userID,
		"pull_requests": shortPRs,
	})
}
