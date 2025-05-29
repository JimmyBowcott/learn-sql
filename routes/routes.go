package routes

import (
	"net/http"
	"fmt"
	"encoding/json"
	"reflect"

	"github.com/JimmyBowcott/learn-sql/database"
)

type SubmitQueryBody struct {
	Query string `json:"query"`
	Level int    `json:"level"`
}

type SubmissionResponse struct {
	Success bool	`json:"success"`
	Result any		`json:"result"`
}

func SubmitQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqBody SubmitQueryBody
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqBody); err != nil {
		http.Error(w, fmt.Sprintf("Failed to read request body: %v", err), http.StatusInternalServerError)
		return
	}

	query := string(reqBody.Query)
	userRes, err := database.ExecuteQuery(query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute query: %v", err), http.StatusBadRequest)
		return
	}

	solution, err := database.GetSolution(reqBody.Level)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get solution: %v", err), http.StatusBadRequest)
		return
	}

	expectedRes, err := database.ExecuteQuery(solution)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get solution: %v", err), http.StatusBadRequest)
		return
	}

	success := reflect.DeepEqual(userRes, expectedRes)

	json.NewEncoder(w).Encode(SubmissionResponse{Success: success, Result: userRes})
}

func GetLevels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	res, err := database.GetLevels()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get levels: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	json.NewEncoder(w).Encode(res)
}
