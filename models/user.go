package models

import "time"

type User struct {
	Id           string
	Username     string
	Email        string
	FirstName    string
	LastName     string
	PasswordHash string
	PasswordSalt string
	CreatedAt    time.Time
	IsActive     bool
}
