package postgresdataaccess

import (
	"database/sql"
	"fmt"
	"todo/models"
)

func (da *PostgresTodoDataAccess) RegisterUser(userName string, email string, firstName string, lastName string, hashedPassword string, salt string) (*models.User, error) {
	var user models.User

	query := `INSERT INTO "user" (user_name, email, first_name, last_name, password_hash, password_salt) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := da.db.Exec(query, userName, email, firstName, lastName, hashedPassword, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	return &user, nil
}

func (da *PostgresTodoDataAccess) GetUserByUsername(username string) (*models.User, error) {
	var user models.User

	query := `SELECT id, user_name, email, first_name, last_name, password_hash, password_salt, created_at, is_active 
              FROM "user" WHERE user_name = $1`

	err := da.db.QueryRow(query, username).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.PasswordHash,
		&user.PasswordSalt,
		&user.CreatedAt,
		&user.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found with username: %s", username)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}
