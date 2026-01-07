package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/SKD-fastcampus/bot-management/services/bot-mgmt-server/internal/domain"
	"github.com/SKD-fastcampus/bot-management/services/bot-mgmt-server/internal/usecase"
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

func (h *TaskHandler) CreateTask(c echo.Context) error {
	var req CreateTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if req.URL == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URL is required"})
	}

	task, err := h.usecase.CreateTask(c.Request().Context(), req.URL, req.RequestUUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusAccepted, task)
}

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
