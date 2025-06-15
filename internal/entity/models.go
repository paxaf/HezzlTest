package entity

import "time"

type Goods struct {
	Id          int
	ProjectId   int
	Name        string
	Description string
	Priority    int
	Removed     bool
	Created_at  time.Time
}
