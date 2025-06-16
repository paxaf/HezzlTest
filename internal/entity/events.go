package entity

import (
	"encoding/json"
	"time"
)

type EventAction string

const (
	Create EventAction = "create"
	Update EventAction = "update"
	Delete EventAction = "delete"
)

type BaseEvent struct {
	Action    EventAction `json:"action"`
	Entity    string      `json:"entity"`
	EntityID  int         `json:"entity_id"`
	ProjectID int         `json:"project_id,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

type Event struct {
	BaseEvent
	Payload interface{} `json:"payload,omitempty"`
}

func (e Event) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

type EventPayload interface {
	ToPayload() interface{}
}

type ProjectEventPayload struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (p Project) ToPayload() interface{} {
	return ProjectEventPayload{
		Name:      p.Name,
		CreatedAt: p.CreatedAt,
	}
}

func NewProjectEvent(action EventAction, project Project) Event {
	return Event{
		BaseEvent: BaseEvent{
			Action:    action,
			Entity:    "project",
			EntityID:  project.Id,
			Timestamp: time.Now().UTC(),
		},
		Payload: project.ToPayload(),
	}
}

type GoodEventPayload struct {
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	CreatedAt   time.Time `json:"created_at"`
}

func (g Goods) ToPayload() interface{} {
	return GoodEventPayload{
		Name:        g.Name,
		Description: &g.Description,
		Priority:    g.Priority,
		Removed:     g.Removed,
		CreatedAt:   g.CreatedAt,
	}
}

func NewGoodEvent(action EventAction, goods Goods) Event {
	return Event{
		BaseEvent: BaseEvent{
			Action:    action,
			Entity:    "good",
			EntityID:  goods.Id,
			ProjectID: goods.ProjectId,
			Timestamp: time.Now().UTC(),
		},
		Payload: goods.ToPayload(),
	}
}
