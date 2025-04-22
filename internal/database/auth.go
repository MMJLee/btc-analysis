package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (repo DBPool) GetUser(c context.Context, username string) (User, error) {
	query := `
		SELECT username FROM auth 
		WHERE username = $1 
	`

	rows, _ := repo.Pool.Query(c, query, username)
	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		return user, fmt.Errorf("GetCandles-%w", err)
	}
	return user, nil
}

func (repo DBPool) CreateUser(c context.Context, username, password string) error {
	query := `
		INSERT INTO auth (username, password) 
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING;
	`
	_, err := repo.Pool.Exec(c, query)
	if err != nil {
		return fmt.Errorf("CreateUser-%w", err)
	}
	return nil
}
