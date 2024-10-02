package services

import "cmd/tambola/internals/repositories"

type authService struct {
	authRepo repositories.AuthRepository
}

type AuthService interface {
	Login(string, string, string) ([]byte, error)
	Signup(string) ([]byte, error)
}

func NewAuthService(authRepo repositories.AuthRepository) AuthService {
	return &authService{
		authRepo: authRepo,
	}
}

func (as *authService) Login(phone string, otp string, firebaseToken string) ([]byte, error) {

}

func (as *authService) Signup(signupToken string) ([]byte, error) {

}

func firebaseVerify() {

}

func generateSignupToken() {

}

func verifySignupToken() {

}

func generateUserToken() {

}
