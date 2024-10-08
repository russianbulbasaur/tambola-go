package handlers

import (
	"cmd/tambola/internals/services"
	"encoding/json"
	"log"
	"net/http"
)

type authHandler struct {
	authService services.AuthService
}

type AuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
}

func NewAuthHandler(authService services.AuthService) AuthHandler {
	return &authHandler{
		authService: authService,
	}
}

func (ah *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		//while first login
		Otp           string `json:"otp"`
		FirebaseToken string `json:"firebase_token"`

		//common
		Phone string `json:"phone"`

		//while signup
		SignupToken string `json:"signup_token"`
		Name        string `json:"name"`
	}
	var req loginRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		log.Println(err)
		return
	}
	response := []byte("")
	if req.SignupToken == "" {
		response, err = ah.authService.Login(req.Phone, req.Otp, req.FirebaseToken)
	} else {
		response, err = ah.authService.Signup(req.Phone, req.SignupToken, req.Name)
	}
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}
