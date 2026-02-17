package repo

import (
	"context"
	"time"
)



type User struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	FullName    string    `json:"full_name"`
	DisplayName string    `json:"display_name,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserSecret struct {
	ID    string `json:"id"`
	Value string `json:"-"` // Never expose secret in JSON
}



type UserRepository interface {
	Create(ctx context.Context, user *User, secret *UserSecret) error
	GetByID(ctx context.Context, id string) (*User, error)

}
