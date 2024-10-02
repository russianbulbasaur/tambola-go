package services

import "cmd/tambola/internals/repositories"

type userService struct {
	userRepo repositories.UserRepository
}

type UserService interface {
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}
