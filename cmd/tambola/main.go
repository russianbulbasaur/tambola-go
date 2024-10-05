/*
The http requests are meant to send data in application/json format.
The form data is used by the middleware to decode tokens and write user objects
for the chained functions.
*/
package main

import (
	db2 "cmd/tambola/db"
	"cmd/tambola/internals/handlers"
	"cmd/tambola/internals/middlewares"
	"cmd/tambola/internals/repositories"
	"cmd/tambola/internals/services"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"log"
	"net/http"
)

var (
	gameHandler handlers.GameHandler
	userHandler handlers.UserHandler
	authHandler handlers.AuthHandler
)

func main() {
	loadEnv()
	db := db2.InitDB()
	userRepository, authRepository := initRepositories(db)
	gameService, userService, authService := initServices(userRepository, authRepository)
	initHandlers(userService, gameService, authService)
	appRouter := initRouter()
	startServer(appRouter)
}

func loadEnv() {
	log.Println("Loading env")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}
}

func initRepositories(db *sql.DB) (repositories.UserRepository, repositories.AuthRepository) {
	userRepo := repositories.NewUserRepository(db)
	authRepo := repositories.NewAuthRepository(db)
	return userRepo, authRepo
}

func initServices(userRepo repositories.UserRepository,
	authRepo repositories.AuthRepository) (services.GameService, services.UserService,
	services.AuthService) {
	gameService := services.NewGameService()
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(authRepo)
	return gameService, userService, authService
}

func initHandlers(userService services.UserService, gameService services.GameService,
	authService services.AuthService) {
	gameHandler = handlers.NewGameHandler(gameService)
	userHandler = handlers.NewUserHandler(userService)
	authHandler = handlers.NewAuthHandler(authService)
}

func initRouter() *http.ServeMux {
	appRouter := http.NewServeMux()
	//non protected routes
	appRouter.HandleFunc("POST /login", authHandler.Login)

	//protected routes
	appRouter.Handle("/game/", http.StripPrefix("/game", gameRouter()))
	appRouter.Handle("/user/", http.StripPrefix("/user", userRouter()))
	return appRouter
}

func gameRouter() http.Handler {
	gameRouter := http.NewServeMux()
	gameRouter.HandleFunc("/create", gameHandler.CreateGame)
	gameRouter.HandleFunc("/join", gameHandler.JoinGame)
	return middlewares.Protect(gameRouter)
}

func userRouter() http.Handler {
	userRouter := http.NewServeMux()
	return middlewares.Protect(userRouter)
}

func startServer(appRouter http.Handler) {
	//port := os.Getenv("PORT")
	port := "8000"
	log.Printf("Starting server at %s", port)
	appRouter = cors.Default().Handler(appRouter)
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), appRouter)
	if err != nil {
		log.Fatalln(err)
	}
}
