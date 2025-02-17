package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service UserService
}

func NewUserHandler(s UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) AddUserRoutes(r *gin.RouterGroup) {

	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/", h.CreateUser)
		userRoutes.PUT("/opt-out/:id", h.Optout)
	}
}

func validateUser(u CreateUserInput) error {
	if u.Name == "" {
		return fmt.Errorf("name missing")
	}
	if u.Email == "" {
		return fmt.Errorf("email missing")
	}
	if u.PhoneNumber == "" {
		return fmt.Errorf("phone_number missing")
	}
	if u.Webhook == "" {
		return fmt.Errorf("webhook missing")
	}
	return nil
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var u CreateUserInput

	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validateUser(u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := AppUser{Name: u.Name, Email: u.Email, PhoneNumber: &u.PhoneNumber, Webhook: &u.Webhook}

	id, err := h.service.CreateUser(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, &CreateUserReturn{id})
}

func (h *UserHandler) Optout(c *gin.Context) {
	var userId int64

	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.OptOut(c, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "User opted-out")
}

type CreateUserInput struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Webhook     string `json:"webhook"`
}

type CreateUserReturn struct {
	Id int64 `json:"id"`
}
