package scheduled_notification

import (
	"fmt"
	"net/http"
	"strings"
	"time"

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
	var n notificationInput
	if err := c.ShouldBindJSON(&n); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validateNotification(n); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	scheduledNotification := ScheduledNotification{Date: n.Date, CityName: n.CityName, UserId: n.UserId, NotificationType: n.NotificationType}

	err := h.service.CreateScheduledNotification(c.Request.Context(), &scheduledNotification)
	if err != nil {
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Used not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification Scheduled",
	})
}

func validateNotification(n notificationInput) error {
	if n.CityName == "" {
		return fmt.Errorf("city_name missing")
	}
	if n.UserId == 0 {
		return fmt.Errorf("user_id missing")
	}
	if n.NotificationType == "" {
		return fmt.Errorf("notification_type missing")
	}
	return nil
}

type notificationInput struct {
	Date             time.Time        `json:"date"`
	CityName         string           `json:"city_name"`
	UserId           int64            `json:"user_id"`
	NotificationType NotificationType `json:"notification_type"`
}
