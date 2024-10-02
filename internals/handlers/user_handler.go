package handlers

import (
	"cmd/tambola/internals/services"
)

type userHandler struct {
	userService services.UserService
}

type UserHandler interface {
}

func NewUserHandler(userService services.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}
