package main

import (
	"cmd/tambola/internals/handlers"
	"cmd/tambola/internals/services"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

var (
	gameHandler handlers.GameHandler
)

func main() {
	//loadEnv()
	gameService := initServices()
	initHandlers(gameService)
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

func initServices() services.GameService {
	gameService := services.NewGameService()
	return gameService
}

func initHandlers(gameService services.GameService) {
	gameHandler = handlers.NewGameHandler(gameService)
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
