package handler

type setIsActiveRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	IsActive bool   `json:"is_active"`
}

type createPRRequest struct {
	ID     string `json:"pull_request_id" binding:"required"`
	Name   string `json:"pull_request_name" binding:"required"`
	Author string `json:"author_id" binding:"required"`
}

type mergePRRequest struct {
	ID string `json:"pull_request_id" binding:"required"`
}

type reassignRequest struct {
	PRID        string `json:"pull_request_id" binding:"required"`
	OldReviewer string `json:"old_user_id" binding:"required"`
}
