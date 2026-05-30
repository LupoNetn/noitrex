package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Envelope struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
}

type Meta struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
	Total   int `json:"total"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Envelope{
		Success: true,
		Data:    data,
	})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Envelope{
		Success: true,
		Data:    data,
	})
}

func BadRequest(c *gin.Context, err string) {
	c.JSON(http.StatusBadRequest, Envelope{
		Success: false,
		Error:   err,
	})
}

func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Envelope{
		Success: false,
		Error:   "unauthorized",
	})
}

func Conflict(c *gin.Context, err string) {
	c.JSON(http.StatusConflict, Envelope{
		Success: false,
		Error:   err,
	})
}

func NotFound(c *gin.Context, resource string) {
	c.JSON(http.StatusNotFound, Envelope{
		Success: false,
		Error:   resource + " not found",
	})
}

func InternalError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, Envelope{
		Success: false,
		Error:   "something went wrong",
	})
}

func Paginated(c *gin.Context, data any, meta Meta) {
	c.JSON(http.StatusOK, Envelope{
		Success: true,
		Data:    data,
		Meta:    &meta,
	})
}
