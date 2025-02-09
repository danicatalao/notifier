package user

import (
	"fmt"
	"net/http"

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
		userRoutes.GET("/", h.CreateUser)
		userRoutes.GET("/:id", h.CreateUser)
		userRoutes.PUT("/", h.CreateUser)
		userRoutes.DELETE("/", h.CreateUser)
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user AppUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	fmt.Printf("%+v\n", user)

	id, err := h.service.CreateUser(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &CreateUserReturn{id})
}

type CreateUserReturn struct {
	Id int64 `json:"id"`
}
