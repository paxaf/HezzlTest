package entity

import "time"

type Goods struct {
	Id          int       `json:"id"`
	ProjectId   int       `json:"project_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	Created_at  time.Time `json:"created_at"`
}

type GoodsResponse struct {
	Goods []Goods `json:"goods"`
}
