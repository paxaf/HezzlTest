package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/paxaf/HezzlTest/internal/entity"
	"github.com/paxaf/HezzlTest/internal/logger"
)

type CreateProject struct {
	Name string `json:"name" binding:"required"`
}

type UpdateProject struct {
	Id   int    `json:"id" binding:"required,gt=0"`
	Name string `json:"name" binding:"required"`
}

func (h *handler) GetProject(c *gin.Context) {
	ctx := c.Request.Context()
	key := c.Request.URL.String()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "Bad request"})
		return
	}
	output, err := h.service.GetProject(ctx, key, id)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			c.JSON(http.StatusNotFound, errorResponse{Error: "Not found"})
			return
		}
		logger.Error("getproject error", err)
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "Internal error"})
		return
	}
	c.JSON(http.StatusOK, output)
}

func (h *handler) GetProjects(c *gin.Context) {
	ctx := c.Request.Context()
	key := c.Request.URL.String()
	output, err := h.service.GetProjects(ctx, key)
	if err != nil {
		logger.Error("getprojects error", err)
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "Internal error"})
		return
	}
	if output == nil {
		c.JSON(http.StatusOK, map[string]interface{}{})
		return
	}
	c.JSON(http.StatusOK, output)
}

func (h *handler) CreateProject(c *gin.Context) {
	ctx := c.Request.Context()
	var req CreateProject
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("JSON binding error", err)
		c.JSON(http.StatusBadRequest, errorResponse{
			Error: "Invalid request format: " + err.Error(),
		})
		return
	}
	input := entity.Project{
		Name: req.Name,
	}
	err := h.service.AddProject(ctx, &input)
	if err != nil {
		logger.Error("createproject error", err)
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "Internal error"})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *handler) UpdateProject(c *gin.Context) {
	ctx := c.Request.Context()
	var req UpdateProject
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "Bad request"})
		return
	}
	input := entity.Project{
		Id:   req.Id,
		Name: req.Name,
	}

	err := h.service.UpdateProject(ctx, &input)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			c.JSON(http.StatusNotFound, errorResponse{Error: "Not found"})
			return
		}
		logger.Error("update proj error", err)
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "Internal error"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *handler) DeleteProject(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "Bad request"})
		return
	}
	err = h.service.DeleteProject(ctx, id)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			c.JSON(http.StatusNotFound, errorResponse{Error: "Not found"})
			return
		}
		logger.Error("deleteproject error", err)
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "Internal error"})
		return
	}
	c.Status(http.StatusNoContent)
}
