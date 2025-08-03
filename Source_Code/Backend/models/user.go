package models

import (
	"SimpleTaskManager/config"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string
}

type SafeUser struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

// Sign Up
func SignUp(user *User) (*SafeUser, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)

	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, email, username`
	var newAccount SafeUser
	row := config.DB.QueryRow(query, user.Username, user.Email, user.Password)
	err = row.Scan(&newAccount.ID, &newAccount.Email, &newAccount.Username)

	if err != nil {
		return nil, err
	}

	return &newAccount, err
}

// Sign In
func SignIn(user *User) (*User, error) {
	var existUser User

	query := `SELECT id, email, password  FROM users WHERE email = $1`

	row := config.DB.QueryRow(query, user.Email)

	err := row.Scan(&existUser.ID, &existUser.Email, &existUser.Password)

	if err != nil {
		return nil, err
	}

	return &existUser, err
}

// Check if account exist
func GetUser(id uuid.UUID) (*SafeUser, error) {
	var user SafeUser
	query := `SELECT id, username, email FROM users WHERE id = $1`
	row := config.DB.QueryRow(query, id)

	err := row.Scan(&user.ID, &user.Username, &user.Email)

	if err != nil {
		return nil, err
	}

	return &user, err
}

// Get All User
func GetAllUser() ([]SafeUser, error) {
	var allUser []SafeUser
	query := `SELECT id, username, email FROM users`
	row, err := config.DB.Query(query)

	if err != nil {
		return nil, err
	}
	defer row.Close()

	for row.Next() {
		var user SafeUser
		err := row.Scan(&user.ID, &user.Username, &user.Email)

		if err != nil {
			return nil, err
		}

		allUser = append(allUser, user)
	}

	return allUser, nil

}

// Delete account

// Update Profile
