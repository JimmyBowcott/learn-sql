package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/JimmyBowcott/learn-sql/routes"
	"github.com/joho/godotenv"
)

func init() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/exec", routes.PostExec)
	mux.HandleFunc("/levels", routes.GetLevels)
	err := http.ListenAndServe(":3456", mux)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
