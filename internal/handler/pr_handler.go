package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) registerPRRoutes(rg *gin.RouterGroup) {
	rg.POST("/create", h.createPR)
	rg.POST("/merge", h.mergePR)
	rg.POST("/reassign", h.reassignReviewer)
}

func (h *Handler) createPR(c *gin.Context) {
	var input createPRRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	pr, err := h.prService.CreatePullRequest(c.Request.Context(), input.ID, input.Name, input.Author)
	if err != nil {
		h.errorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"pr": pr})
}

func (h *Handler) mergePR(c *gin.Context) {
	var input mergePRRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	pr, err := h.prService.MergePullRequest(c.Request.Context(), input.ID)
	if err != nil {
		h.errorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"pr": pr})
}

func (h *Handler) reassignReviewer(c *gin.Context) {
	var input reassignRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	pr, newReviewerID, err := h.prService.ReassignReviewer(c.Request.Context(), input.PRID, input.OldReviewer)
	if err != nil {
		h.errorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pr":          pr,
		"replaced_by": newReviewerID,
	})
}
