package main

import (
	db2 "cmd/tambola/db"
	"cmd/tambola/internals/handlers"
	"cmd/tambola/internals/repositories"
	"cmd/tambola/internals/services"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

var (
	gameHandler handlers.GameHandler
	userHandler handlers.UserHandler
)

func main() {
	//loadEnv()
	db := db2.InitDB()
	userRepository := initRepositories(db)
	gameService, userService := initServices(userRepository)
	initHandlers(userService, gameService)
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

func initRepositories(db *sql.DB) repositories.UserRepository {
	userRepo := repositories.NewUserRepository(db)
	return userRepo
}

func initServices(userRepo repositories.UserRepository) (services.GameService, services.UserService) {
	gameService := services.NewGameService()
	userService := services.NewUserService(userRepo)
	return gameService, userService
}

func initHandlers(userService services.UserService, gameService services.GameService) {
	gameHandler = handlers.NewGameHandler(gameService)
	userHandler = handlers.NewUserHandler(userService)
}

func initRouter() *http.ServeMux {
	appRouter := http.NewServeMux()
	appRouter.Handle("/game/", http.StripPrefix("/game", gameRouter()))
	return appRouter
}

func gameRouter() *http.ServeMux {
	gameRouter := http.NewServeMux()
	gameRouter.HandleFunc("/create", gameHandler.CreateGame)
	gameRouter.HandleFunc("/join", gameHandler.JoinGame)
	return gameRouter
}

func startServer(appRouter *http.ServeMux) {
	//port := os.Getenv("PORT")
	port := "8000"
	log.Printf("Starting server at %s", port)
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), appRouter)
	if err != nil {
		log.Fatalln(err)
	}
}
