package main

import (
	"net/http"
	"fmt"
	"errors"
	"os"

	"github.com/JimmyBowcott/learn-sql/routes"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/exec", routes.PostExec)
	err := http.ListenAndServe(":3456", mux)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
