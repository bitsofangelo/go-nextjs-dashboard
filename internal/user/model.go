package user

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Name     string
	Email    string
	Password string
}

func (u User) UserID() any {
	return u.ID
}
