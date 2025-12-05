package storage

import (
	"auth-service/internal/model"
	"context"
	"database/sql"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	user := &model.User{}

	query := `
			SELECT id, username, email, role, created_at, updated_at
			FROM auth.users
			WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.Id, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepo) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	user := &model.User{}

	query := `
		SELECT id, username, email, role, password_hash, created_at, updated_at
		FROM auth.users
		WHERE username = $1 OR email = $1
	`

	err := r.db.QueryRowContext(ctx, query, login).Scan(
		&user.Id, &user.Username, &user.Email, &user.Role, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepo) CreateUser(ctx context.Context, User *model.User) (*model.User, error) {
	User.Role = model.RoleUser

	userResp := &model.User{}

	query := `
		INSERT INTO auth.users (username, email, password_hash, role)
		VALUES	($1, $2, $3, $4)
		RETURNING id, username, email, role, created_at, updated_at 
	`

	err := r.db.QueryRowContext(ctx, query, User.Username, User.Email, User.PasswordHash, User.Role).Scan(
		&userResp.Id, &userResp.Username, &userResp.Email, &userResp.Role, &userResp.CreatedAt, &userResp.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return userResp, nil
}
