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
	EntityID  string      `json:"entity_id"`
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
