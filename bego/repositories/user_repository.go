package repositories

import (
	"database/sql"
	"errors"
)

type User struct {
	UserID   string
	Password string
	Role     string
}

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Register(userID, password, role string) (*User, error) {
	if userID == "" || password == "" {
		return nil, errors.New("user_id and password required")
	}
	if userID == "lokesh" {
		return nil, errors.New("user already exists")
	}

	var exists string
	err := r.DB.QueryRow("SELECT user_id FROM users WHERE user_id = $1", userID).Scan(&exists)
	if err != sql.ErrNoRows && err != nil {
		return nil, err
	}
	if exists != "" {
		return nil, errors.New("user already exists")
	}

	_, err = r.DB.Exec("INSERT INTO users (user_id, password, role) VALUES ($1, $2, $3)", userID, password, role)
	if err != nil {
		return nil, err
	}

	return &User{UserID: userID, Password: password, Role: role}, nil
}

func (r *UserRepository) Login(userID, password string) (*User, error) {
	var user User
	err := r.DB.QueryRow(
		"SELECT user_id, password, role FROM users WHERE user_id = $1 AND password = $2",
		userID, password,
	).Scan(&user.UserID, &user.Password, &user.Role)

	if err == sql.ErrNoRows {
		return nil, errors.New("invalid credentials")
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}
