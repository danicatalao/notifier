package scheduled_notification

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ScheduledNotificationHandler struct {
	service ScheduledNotificationService
}

func NewScheduledNotificationHandler(s ScheduledNotificationService) *ScheduledNotificationHandler {
	return &ScheduledNotificationHandler{service: s}
}

func (sn *ScheduledNotificationHandler) AddNotificationRoutes(r *gin.RouterGroup) {

	userRoutes := r.Group("/notification")
	{
		userRoutes.POST("/", sn.CreateScheduledNotification)
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
