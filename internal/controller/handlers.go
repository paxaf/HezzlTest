package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/paxaf/HezzlTest/internal/entity"
	"github.com/paxaf/HezzlTest/internal/usecase"
)

type handler struct {
	service usecase.Usecase
}

func New(service usecase.Usecase) *handler {
	return &handler{
		service: service,
	}
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *handler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	key := c.Request.URL.String()
	output, err := h.service.GetAllItems(ctx, key)
	if err != nil {
		// log
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "Internal error"})
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
		// log
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
		//log
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "Internal error"})
		return
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
		// log
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "Internal error"})
		return
	}
	c.JSON(http.StatusOK, output)
}

func (h *handler) CreateItem(c *gin.Context) {
	ctx := c.Request.Context()
	var input *entity.Goods
	projectIdStr := c.Request.FormValue("project_id")
	input.Name = c.Request.FormValue("name")
	input.Description = c.Request.FormValue("description")
	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "Bad request"})
		return
	}
	input.ProjectId = projectId
	err = h.service.CreateItem(ctx, input)
	if err != nil {
		// log
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "Internal error"})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *handler) UpdateItem(c *gin.Context) {
	ctx := c.Request.Context()
	var input *entity.Goods
	input.Name = c.Request.FormValue("name")
	input.Description = c.Request.FormValue("description")
	priorityStr := c.Request.FormValue("priority")
	removedStr := c.Request.FormValue("removed")
	priority, err := strconv.Atoi(priorityStr)
	if err != nil || priority < 1 {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "Bad request"})
		return
	}
	removed, err := strconv.ParseBool(removedStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "Bad request"})
		return
	}
	input.Priority = priority
	input.Removed = removed
	err = h.service.UpdateItem(ctx, input)
	if err != nil {
		if errors.As(err, entity.ErrNotFound) {
			c.JSON(http.StatusNotFound, errorResponse{Error: "Not found"})
			return
		}
		//log
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
		if errors.As(err, entity.ErrNotFound) {
			c.JSON(http.StatusNotFound, errorResponse{Error: "Not found"})
			return
		}
		//log
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "Internal error"})
		return
	}
	c.Status(http.StatusNoContent)
}
