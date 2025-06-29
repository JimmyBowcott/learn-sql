package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/JimmyBowcott/learn-sql/routes"
	"github.com/JimmyBowcott/learn-sql/middleware"
	"github.com/joho/godotenv"
)

func Secure(h http.HandlerFunc) http.HandlerFunc {
	return middleware.CorsMiddleware(middleware.RateLimitMiddleware(h))
}

func init() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/submit", Secure(routes.SubmitQuery))
	mux.HandleFunc("/levels", Secure(routes.GetLevels))
	mux.HandleFunc("/signin", Secure(routes.Login))
	mux.HandleFunc("/signup", Secure(routes.SignUp))
	err := http.ListenAndServe(":3456", mux)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
