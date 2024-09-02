package database

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"tincho.dev/rest-ws/models"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(url string) (*Postgres, error) {
	db, err := sql.Open("postgres", url)

	if err != nil {
		return nil, err
	}

	return &Postgres{db}, nil
}

func (p *Postgres) CreateUser(ctx context.Context, user *models.User) error {
	_, err := p.db.ExecContext(ctx, "INSERT INTO users (email, password) VALUES ($1, $2)", user.Email, user.Password)

	return err
}

func (p *Postgres) FindAllUsers(ctx context.Context) ([]models.User, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT id, email, password FROM users")

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

func (p *Postgres) Close() error {
	return p.db.Close()
}
