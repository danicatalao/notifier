package scheduled_notification

import (
	"fmt"
	"net/http"

	l "github.com/danicatalao/notifier/internal/logger"
	"github.com/gin-gonic/gin"
)

type ScheduledNotificationHandler struct {
	service ScheduledNotificationService
	log     l.Logger
}

func NewScheduledNotificationHandler(s ScheduledNotificationService, l l.Logger) *ScheduledNotificationHandler {
	return &ScheduledNotificationHandler{service: s, log: l}
}

func (h *ScheduledNotificationHandler) AddNotificationRoutes(r *gin.RouterGroup) {

	userRoutes := r.Group("/notification")
	{
		userRoutes.POST("/", h.CreateScheduledNotification)
	}
}

func (h *ScheduledNotificationHandler) CreateScheduledNotification(c *gin.Context) {
	var scheduledNotification ScheduledNotification
	if err := c.ShouldBindJSON(&scheduledNotification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	fmt.Printf("%+v\n", scheduledNotification)

	err := h.service.CreateScheduledNotification(c.Request.Context(), &scheduledNotification)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification Scheduled",
	})
}
