package postgresdataaccess

import (
	"fmt"
	"todo/models"
)

func (da *PostgresTodoDataAccess) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := da.db.QueryRow(`SELECT id, username, passwordhash, passwordsalt FROM user WHERE username = '$1';`, username).Scan(&user.Id, &user.Username, &user.PasswordHash, &user.PasswordSalt)
	if err != nil {
		return &user, fmt.Errorf("Unable to find the user with username: %s: %w", username, err)
	}

	return nil, nil
}
