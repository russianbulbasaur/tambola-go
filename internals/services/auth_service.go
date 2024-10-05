package services

import (
	"bytes"
	"cmd/tambola/internals/repositories"
	"cmd/tambola/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type authService struct {
	authRepo repositories.AuthRepository
}

type AuthService interface {
	Login(string, string, string) ([]byte, error)
	Signup(string, string, string) ([]byte, error)
}

func NewAuthService(authRepo repositories.AuthRepository) AuthService {
	return &authService{
		authRepo: authRepo,
	}
}

func (as *authService) Login(phone string, otp string, firebaseToken string) ([]byte, error) {
	if firebaseVerify(otp, firebaseToken, phone) {
		user, err := as.authRepo.FindUser(phone)
		if err != nil {
			return nil, err
		}
		if user.Id == 0 {
			//signup
			signupToken := generateSignupToken(phone)
			user = &models.User{
				Phone: phone,
				Token: signupToken,
			}
		} else {
			userToken := generateUserToken(user)
			user.Token = userToken
		}
		return user.EncodeToJson(), nil
	}
	return nil, errors.New("firebase token invalid")
}

func (as *authService) Signup(phone string, signupToken string, name string) ([]byte, error) {
	if verifySignupToken(signupToken, phone) {
		user, err := as.authRepo.Signup(name, phone)
		if err != nil {
			log.Fatalln(err)
		}
		userToken := generateUserToken(user)
		user.Token = userToken
		return user.EncodeToJson(), nil
	}
	return nil, errors.New("signup token invalid")
}

func generateUserToken(user *models.User) string {
	data := jwt.MapClaims{
		"sub": string(user.EncodeToJson()),
	}
	return generateToken(data)
}

func generateSignupToken(phone string) string {
	data := jwt.MapClaims{
		"sub": phone,
		"exp": time.Now().Add(time.Minute * 5).Unix(),
	}
	return generateToken(data)
}

func generateToken(data jwt.MapClaims) string {
	key := os.Getenv("JWT_SECRET")
	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	token, err := tokenWithClaims.SignedString([]byte(key))
	if err != nil {
		log.Fatalln(err)
	}
	return token
}

func verifySignupToken(token string, phone string) bool {
	data, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		key := []byte(os.Getenv("JWT_SECRET"))
		return key, nil
	})
	if err != nil || !data.Valid {
		log.Println(err)
		return false
	}
	sub, err := data.Claims.GetSubject()
	if err != nil || sub != phone {
		log.Println(err)
		return false
	}
	return true
}

func firebaseVerify(otp string, firebaseToken string, phone string) bool {

	type FirebaseResponse struct {
		Phone string `json:"phoneNumber"`
	}

	apiKey := os.Getenv("FIREBASE_API_KEY")
	body, _ := json.Marshal(map[string]string{
		"sessionInfo": firebaseToken,
		"code":        otp,
	})
	urlString := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPhoneNumber?key=%s", apiKey)
	client := &http.Client{}
	response, err := client.Post(urlString, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
		return false
	}
	defer response.Body.Close()
	body, _ = io.ReadAll(response.Body)
	var firebaseRes FirebaseResponse
	err = json.Unmarshal(body, &firebaseRes)
	if err != nil {
		log.Println(err)
		return false
	}
	if firebaseRes.Phone == phone {
		return true
	}
	return false
}
