package repository

import (
	"context"
	"tiket-event/models/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (repository *userRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	querySql := "INSERT INTO users (name, email, password) VALUES($1, $2, $3) RETURNING id, role, created_at"
	err := repository.db.QueryRow(ctx, querySql, user.Name, user.Email, user.Password).Scan(&user.Id, &user.Role, &user.CreateAt)
	return user, err
}

func (repostory *userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	querySql := `SELECT id, name, email, password, role, created_at FROM users WHERE email = $1`
	err := repostory.db.QueryRow(ctx, querySql, email).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreateAt)
	return user, err
}
