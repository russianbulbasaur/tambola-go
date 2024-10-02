package handlers

import (
	"cmd/tambola/internals/services"
	"net/http"
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

func (uh *userHandler) login(w http.ResponseWriter, r *http.Request) {

}
