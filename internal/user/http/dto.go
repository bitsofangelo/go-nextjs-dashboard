package http

import (
	"github.com/google/uuid"

	"go-nextjs-dashboard/internal/user"
)

type Response struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

func toResponse(u *user.User) Response {
	return Response{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}
