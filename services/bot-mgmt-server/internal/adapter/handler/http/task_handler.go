package http

import (
	"net/http"
	"net/url"

	"github.com/SKD-fastcampus/bot-management/services/bot-mgmt-server/internal/domain"
	"github.com/SKD-fastcampus/bot-management/services/bot-mgmt-server/internal/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TaskHandler struct {
	usecase usecase.TaskUsecase
}

func NewTaskHandler(u usecase.TaskUsecase) *TaskHandler {
	return &TaskHandler{usecase: u}
}

// RegisterRoutes registers the task routes with the echo group
func (h *TaskHandler) RegisterRoutes(g *echo.Group) {
	g.POST("/analyze", h.CreateTask)
	g.GET("/status/:id", h.GetStatus)
	g.POST("/webhook", h.HandleWebhook)
}

type CreateTaskRequest struct {
	URL         string `json:"url"`
	RequestUUID string `json:"request_uuid"` // Can be JWT or simple UUID
}

// CreateTask godoc

// @Summary Create a new analysis task
// @Description Initiates a new smishing analysis task for a given URL
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body CreateTaskRequest true "Create Task Request"
// @Success 202 {object} domain.AnalysisTask
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /analyze [post]
func (h *TaskHandler) CreateTask(c echo.Context) error {
	var req CreateTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if req.URL == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URL is required"})
	}

	// Validate URL format
	parsedURL, err := url.Parse(req.URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid URL format"})
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URL must use http or https scheme"})
	}

	task, err := h.usecase.CreateTask(c.Request().Context(), req.URL, req.RequestUUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusAccepted, task)
}

// GetStatus godoc
// @Summary Get task status
// @Description Retrieve the current status of an analysis task
// @Tags tasks
// @Produce json
// @Param id path string true "Task ID" format(uuid)
// @Success 200 {object} domain.AnalysisTask
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /status/{id} [get]
func (h *TaskHandler) GetStatus(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid task ID"})
	}

	task, err := h.usecase.GetTaskStatus(c.Request().Context(), id)
	if err != nil {
		// Differentiate between 404 and 500
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, task)
}

type WebhookRequest struct {
	TaskID string            `json:"task_id"` // Matches our internal ID
	Status domain.TaskStatus `json:"status"`
	Result string            `json:"result"`
}

// HandleWebhook godoc

// @Summary Handle webhook update
// @Description Update task status via webhook (Internal use)
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body WebhookRequest true "Webhook Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /webhook [post]
func (h *TaskHandler) HandleWebhook(c echo.Context) error {
	var req WebhookRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid webhook payload"})
	}

	id, err := uuid.Parse(req.TaskID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid task ID"})
	}

	if err := h.usecase.UpdateTaskStatus(c.Request().Context(), id, req.Status, req.Result); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// If webhook reports failure, should we trigger retry immediately?
	// The background worker will pick it up, or we can trigger it explicitly if logic requires.
	// For now, let background worker handle retry on next tick or if we add logic to UpdateTaskStatus.

	return c.JSON(http.StatusOK, map[string]string{"message": "Status updated"})
}
