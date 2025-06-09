package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/JimmyBowcott/learn-sql/auth"
	"github.com/JimmyBowcott/learn-sql/database"
)

type SubmitQueryBody struct {
	Query string `json:"query"`
	Level int    `json:"level"`
}

type SubmissionResponse struct {
	Success bool   `json:"success"`
	Result  any    `json:"result"`
	Token   string `json:"token"`
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
		http.Error(w, fmt.Sprintf("Failed to get solution: %v", err), http.StatusInternalServerError)
		return
	}

	expectedRes, err := database.ExecuteQuery(solution)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get solution: %v", err), http.StatusInternalServerError)
		return
	}

	success := reflect.DeepEqual(userRes, expectedRes)
	if !success {
		json.NewEncoder(w).Encode(SubmissionResponse{Success: success, Result: userRes, Token: ""})
	}

	claims := auth.GetClaims(r)
	if reqBody.Level > 1 && claims.Level < reqBody.Level {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(claims.Username, reqBody.Level+1)

	json.NewEncoder(w).Encode(SubmissionResponse{Success: success, Result: userRes, Token: token})
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

type Credentials struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

type LoginResponse struct {
	Name  string `json:"name"`
	Level int    `json:"level"`
	Token string `json:"token"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqBody Credentials
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqBody); err != nil {
		http.Error(w, fmt.Sprintf("Failed to read request body: %v", err), http.StatusInternalServerError)
		return
	}

	user, err := database.GetUser(reqBody.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to find user: %v", err), http.StatusBadRequest)
		return
	}

	validPassword := database.ValidateUser(user, reqBody.Pass)
	if !validPassword {
		http.Error(w, "Password is incorrect", http.StatusUnauthorized)
		return
	}

	claims := auth.GetClaims(r)
	fmt.Println(user.Level, claims.Level)
	if claims.Level > user.Level {
		user.Level = claims.Level
		database.SetLevel(user.Name, claims.Level)
	}

	json.NewEncoder(w).Encode(LoginResponse{Name: user.Name, Level: user.Level, Token: user.Token})
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqBody Credentials
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqBody); err != nil {
		http.Error(w, fmt.Sprintf("Failed to read request body: %v", err), http.StatusInternalServerError)
		return
	}

	claims := auth.GetClaims(r)
	err := database.CreateUser(reqBody.Name, reqBody.Pass, claims.Level)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
