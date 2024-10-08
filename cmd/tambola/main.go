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
	_ "net/http/pprof"
	"runtime"
	"time"
)

var (
	gameHandler handlers.GameHandler
	userHandler handlers.UserHandler
	authHandler handlers.AuthHandler
)

func main() {
	loadEnv()
	db := db2.InitDB()
	userRepository, authRepository, gameRepository := initRepositories(db)
	gameService, userService, authService := initServices(userRepository, authRepository, gameRepository)
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

func initRepositories(db *sql.DB) (repositories.UserRepository, repositories.AuthRepository,
	repositories.GameRepository) {
	userRepo := repositories.NewUserRepository(db)
	authRepo := repositories.NewAuthRepository(db)
	gameRepo := repositories.NewGameRepository(db)
	return userRepo, authRepo, gameRepo
}

func initServices(userRepo repositories.UserRepository,
	authRepo repositories.AuthRepository, gameRepo repositories.GameRepository) (services.GameService, services.UserService,
	services.AuthService) {
	gameService := services.NewGameService(gameRepo)
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
	//appRouter.Handle("/public/", http.FileServer(http.Dir("/public/")))

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
	go func() {
		for {
			log.Printf("Goroutine count : %d", runtime.NumGoroutine())
			time.Sleep(time.Second * 5)
		}
	}()
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), appRouter)
	if err != nil {
		log.Fatalln(err)
	}
}
