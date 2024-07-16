package handler

import (
	"go-nextjs-dashboard/config"
	"go-nextjs-dashboard/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB          *gorm.DB
	userService *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{DB: config.DB, userService: service.NewUserService()}
}

func (h *UserHandler) GetUserByEmail(c *fiber.Ctx) error {
	email := c.Params("email")

	if err := config.Validate.Var(email, "email"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid email"})
	}

	user, err := h.userService.GetUserByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"message": "User not found"})
		}

		return err
	}

	userResp := struct {
		ID    uuid.UUID `json:"id"`
		Name  string    `json:"name"`
		Email string    `json:"email"`
	}{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	return c.Status(200).JSON(fiber.Map{"data": userResp})
}
