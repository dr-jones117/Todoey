package models

type User struct {
	Id           string
	Username     string
	PasswordHash string
	PasswordSalt string
}
