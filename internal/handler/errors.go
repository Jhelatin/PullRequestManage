package handler

import (
	"PullRequestManage/internal/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) errorResponse(c *gin.Context, err error) {
	h.logger.Printf("ERROR: %v", err) // Логируем ошибку

	errorBody := func(code, message string) gin.H {
		return gin.H{"error": gin.H{"code": code, "message": message}}
	}

	switch {
	case errors.Is(err, service.ErrTeamExists):
		c.JSON(http.StatusBadRequest, errorBody("TEAM_EXISTS", err.Error()))
	case errors.Is(err, service.ErrPRExists):
		c.JSON(http.StatusConflict, errorBody("PR_EXISTS", err.Error()))
	case errors.Is(err, service.ErrPRMerged):
		c.JSON(http.StatusConflict, errorBody("PR_MERGED", err.Error()))
	case errors.Is(err, service.ErrReviewerNotAssigned):
		c.JSON(http.StatusConflict, errorBody("NOT_ASSIGNED", err.Error()))
	case errors.Is(err, service.ErrNoCandidateForReassign):
		c.JSON(http.StatusConflict, errorBody("NO_CANDIDATE", err.Error()))
	case errors.Is(err, service.ErrTeamNotFound),
		errors.Is(err, service.ErrUserNotFound),
		errors.Is(err, service.ErrAuthorNotFound),
		errors.Is(err, service.ErrPRNotFound):
		c.JSON(http.StatusNotFound, errorBody("NOT_FOUND", err.Error()))
	default:
		c.JSON(http.StatusInternalServerError, errorBody("INTERNAL_ERROR", "an unexpected internal error occurred"))
	}

	// todo: переделать
}
