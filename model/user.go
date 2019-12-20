package model

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type UserID string

func newUserID() UserID {
	return UserID(fmt.Sprintf("User:%s", uuid.New().String()))
}

func IsUserID(s string) bool {
	return strings.Contains(s, "User")
}

func NewUser(name string) *User {
	return &User{
		ID:   newUserID(),
		Name: name,
	}
}

type User struct {
	ID   UserID `boom:"id" json:"id"`
	Name string `json:"name"`
}

func (u *User) IsNode() {}

type UserConnection struct {
	Edges    []*UserEdge `json:"edges"`
	PageInfo *PageInfo   `json:"pageInfo"`
}

type UserEdge struct {
	Cursor string `json:"cursor"`
	Node   *User  `json:"node"`
}
