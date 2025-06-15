package controller

import (
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
