package database

import (
	"database/sql"
	"os"
 
	"github.com/JimmyBowcott/learn-sql/auth"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name string
	Pass string
	Level int
	Token string
}

func GetUser(username string) (User, error) {
	connStr := os.Getenv("DB_CONNECTION_STRING_2")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return User{}, err
	}
	defer db.Close()

	var user User
	err = db.QueryRow("SELECT name, pass, level FROM app_user WHERE name = $1;", username).
			Scan(&user.Name, &user.Pass, &user.Level)
	if err != nil {
		return User{}, err
	}

	user.Token, err = auth.GenerateToken(user.Name, user.Level)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func ValidateUser(user User, pass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(pass)) == nil
}

func CreateUser(username string, pass string) error {
	connStr := os.Getenv("DB_CONNECTION_STRING_2")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	hashedPass, err := hashPassword(pass)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO app_user (name, pass, level) VALUES ($1, $2, 1)", username, hashedPass)
	if err != nil {
		return err
	}
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
