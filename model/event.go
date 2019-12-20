package model

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type EventID string

func newEventID() EventID {
	return EventID(fmt.Sprintf("Event:%s", uuid.New().String()))
}

func IsEventID(s string) bool {
	return strings.Contains(s, "Event")
}

func NewEvent(userID UserID, description string) *Event {
	return &Event{
		ID:          newEventID(),
		UserID:      userID,
		Description: description,
	}
}

type Event struct {
	ID          EventID `boom:"id" json:"id"`
	UserID      UserID  `json:"userID"`
	Description string  `json:"description"`
}

func (e *Event) IsNode() {}

type EventConnection struct {
	Edges    []*EventEdge `json:"edges"`
	PageInfo *PageInfo    `json:"pageInfo"`
}

type EventEdge struct {
	Cursor string `json:"cursor"`
	Node   *Event `json:"node"`
}
