package routes

import (
	"net/http"
	"io"
	"fmt"
	"encoding/json"

	"github.com/JimmyBowcott/learn-sql/database"
)

func PostExec(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read request body: %v", err), http.StatusBadRequest)
		return
	}

	query := string(body)
	res, err := database.ExecuteQuery(query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute querry: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	json.NewEncoder(w).Encode(res)
}
