package database

import (
	"context"

	"tincho.dev/rest-ws/models"
)

func (p *Postgres) CreateUser(ctx context.Context, user *models.User) error {
	_, err := p.db.ExecContext(ctx, "INSERT INTO users (email, password) VALUES ($1, $2)", user.Email, user.Password)

	return err
}

func (p *Postgres) ListUsers(ctx context.Context, offset int64, limit int64) ([]models.User, error) {
	query := `
		SELECT id, email, password
		FROM users
		OFFSET $1
		LIMIT $2
	`

	rows, err := p.db.QueryContext(ctx, query, limit*offset, limit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []models.User{}

	for rows.Next() {
		var u models.User

		err := rows.Scan(&u.Id, &u.Email, &u.Password)

		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func (p *Postgres) FindAllUsers(ctx context.Context) ([]models.User, error) {
	query := `
		SELECT id, email, password
		FROM users
	`

	rows, err := p.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []models.User{}

	for rows.Next() {
		var u models.User

		err := rows.Scan(&u.Id, &u.Email, &u.Password)

		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func (p *Postgres) FindUserById(ctx context.Context, id int64) (*models.User, error) {
	row := p.db.QueryRowContext(ctx, "SELECT id, email, password FROM users WHERE id = $1", id)

	var u models.User

	err := row.Scan(&u.Id, &u.Email, &u.Password)

	if err != nil {
		return nil, err
	}

	posts, err := p.FindPostsByUserId(ctx, id)

	if err != nil {
		return nil, err
	}

	u.Posts = posts

	return &u, nil
}

func (p *Postgres) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	row := p.db.QueryRowContext(ctx, "SELECT id, email, password FROM users WHERE email = $1", email)

	var u models.User

	err := row.Scan(&u.Id, &u.Email, &u.Password)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (p *Postgres) UpdateOneUser(ctx context.Context, user *models.User) error {
	_, err := p.db.ExecContext(ctx, "UPDATE users SET email = $1, password = $2 WHERE id = $3", user.Email, user.Password, user.Id)

	return err
}

func (p *Postgres) DeleteOneUser(ctx context.Context, id int64) error {
	_, err := p.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)

	return err
}
