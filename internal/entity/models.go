package entity

import (
	"errors"
	"time"
)

type Project struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Goods struct {
	Id          int       `json:"id"`
	ProjectId   int       `json:"project_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	CreatedAt   time.Time `json:"created_at"`
}

type GoodsResponse struct {
	Goods []Goods `json:"goods"`
}

type ProjectResponse struct {
	Project []Project `json:"project"`
}

var ErrNotFound = errors.New("not found")
