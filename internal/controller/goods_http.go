package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/paxaf/HezzlTest/internal/entity"
	"github.com/paxaf/HezzlTest/internal/logger"
)

type CreateRequest struct {
	ProjectID   int    `json:"project_id" binding:"required,gt=0"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateRequset struct {
	Id          int    `json:"id" binding:"required,gt=0"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Priority    int    `json:"priority" binding:"required,gt=0"`
	Removed     bool   `json:"removed" binding:"required"`
}

func (h *handler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	key := c.Request.URL.String()
	output, err := h.service.GetAllItems(ctx, key)
	if err != nil {
		logger.Error("getall rrror", err)
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "Internal error"})
		return
	}
	if output == nil {
		c.JSON(http.StatusOK, map[string]interface{}{})
		return
	}
	c.JSON(http.StatusOK, output)
}

func (h *handler) GetItem(c *gin.Context) {
	ctx := c.Request.Context()
	key := c.Request.URL.String()
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "Bad request"})
		return
	}
	output, err := h.service.GetItem(ctx, key, id)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			c.JSON(http.StatusNotFound, errorResponse{Error: "Not found"})
			return
		}
		logger.Error("getitem error", err)
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "Internal error"})
		return
	}
	c.JSON(http.StatusOK, output)
}

func (h *handler) GetItemsByName(c *gin.Context) {
	ctx := c.Request.Context()
	key := c.Request.URL.String()
	name := c.Param("name")
	output, err := h.service.GetItemsByName(ctx, key, name)
	if err != nil {
		logger.Error("getitemsbyname error", err)
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "Internal error"})
		return
	}
	if output == nil {
		c.JSON(http.StatusOK, map[string]interface{}{})
	}
	c.JSON(http.StatusOK, output)
}

func (h *handler) GetItemsByProject(c *gin.Context) {
	ctx := c.Request.Context()
	key := c.Request.URL.String()
	projectStr := c.Param("project_id")
	projectId, err := strconv.Atoi(projectStr)
	if err != nil || projectId < 1 {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "Bad request"})
		return
	}
	output, err := h.service.GetItemsByProject(ctx, key, projectId)
	if err != nil {
		logger.Error("getitemsbyproject error", err)
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "Internal error"})
		return
	}
	if output == nil {
		c.JSON(http.StatusOK, map[string]interface{}{})
	}
	c.JSON(http.StatusOK, output)
}

func (h *handler) CreateItem(c *gin.Context) {
	ctx := c.Request.Context()
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("JSON binding error", err)
		c.JSON(http.StatusBadRequest, errorResponse{
			Error: "Invalid request format: " + err.Error(),
		})
		return
	}
	input := entity.Goods{
		ProjectId:   req.ProjectID,
		Description: req.Description,
		Name:        req.Name,
	}
	err := h.service.CreateItem(ctx, &input)
	if err != nil {
		logger.Error("createitem error", err)
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "Internal error"})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *handler) UpdateItem(c *gin.Context) {
	ctx := c.Request.Context()
	var req UpdateRequset
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "Bad request"})
		return
	}
	input := entity.Goods{
		Id:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Removed:     req.Removed,
		Priority:    req.Priority,
	}

	err := h.service.UpdateItem(ctx, &input)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			c.JSON(http.StatusNotFound, errorResponse{Error: "Not found"})
			return
		}
		logger.Error("update error", err)
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "Internal error"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *handler) DeleteItem(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "Bad request"})
		return
	}
	err = h.service.DeleteItem(ctx, id)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			c.JSON(http.StatusNotFound, errorResponse{Error: "Not found"})
			return
		}
		logger.Error("deleteitem error", err)
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "Internal error"})
		return
	}
	c.Status(http.StatusNoContent)
}
