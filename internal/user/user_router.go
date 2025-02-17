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
		userRoutes.POST("/:id/opt-out", h.CreateUser)
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

func (h *UserHandler) Optout(c *gin.Context) {
	var userId int64

	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.service.OptOut(c, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Usuário desativado")
}

type CreateUserReturn struct {
	Id int64 `json:"id"`
}
