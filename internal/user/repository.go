package user

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	//go:embed sql/create_user.sql
	createUserQuery string
	//go:embed sql/get_user_by_id.sql
	getUserByIdQuery string
	//go:embed sql/delete_user_by_id.sql
	deleteUserByIdQuery string
	//go:embed sql/update_user_by_id.sql
	updateUserByIdQuery string
	//go:embed sql/get_user_by_username.sql
	getUserByUsernameQuery string
	//go:embed sql/update_user_by_username.sql
	updateUserByUsernameQuery string
	//go:embed sql/delete_user_by_username.sql
	deleteUserByUsernameQuery string
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) CreateUser(ctx context.Context, username string, email string) (User, error) {
	var user User
	err := r.pool.QueryRow(ctx, createUserQuery, username, email).
		Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.Deleted,
		)

	if err != nil {
		slog.Error("[UserRepository-CreateUser]", "Error", err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// Duplicate key value violates unique constraint
			if pgErr.Code == "23505" {
				return User{}, &UsernameIsTakenError{}
			}
		}

		// Need to check other errors, it might be due to unique username
		return User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, userId string) (User, error) {
	var user User
	err := r.pool.QueryRow(ctx, getUserByIdQuery, userId).
		Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.Deleted,
		)

	if err != nil {
		slog.Error("[UserRepository-GetUserById]", "Error", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, &UserDoesNotExistError{}
		}

		return User{}, errors.New("unknown error when trying to GET user by ID")
	}

	return user, nil
}

func (r *UserRepository) DeleteUserById(ctx context.Context, userId string) (User, error) {
	var user User
	err := r.pool.QueryRow(ctx, deleteUserByIdQuery, userId).
		Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.Deleted,
		)

	if err != nil {
		slog.Error("[UserRepository-DeleteUserById]", "Error", err)

		// Need to check other errors if there is any
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, &UserDoesNotExistError{}
		}

		return User{}, errors.New("unknown error when trying to DELETE user by ID")
	}

	return user, nil
}

func (r *UserRepository) UpdateUserById(ctx context.Context, userId *string, newUsername *string, newEmail *string) (User, error) {
	var setClauses []string
	var args []interface{}
	paramIndex := 1

	if newUsername != nil {
		setClauses = append(setClauses, fmt.Sprintf("username = $%d", paramIndex))
		args = append(args, *newUsername)
		paramIndex++
	}
	if newEmail != nil {
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", paramIndex))
		args = append(args, *newEmail)
		paramIndex++
	}

	setClauses = append(setClauses, "updated_at = NOW()")
	args = append(args, *userId)

	query := fmt.Sprintf(
		updateUserByIdQuery,
		strings.Join(setClauses, ", "),
		paramIndex,
	)

	var user User
	err := r.pool.QueryRow(ctx, query, args...).
		Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.Deleted,
		)

	if err != nil {
		slog.Error("[UserRepository-UpdateUserById]", "Error", err)

		// Need to check other errors if there is any
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, &UserDoesNotExistError{}
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// Duplicate key value violates unique constraint
			if pgErr.Code == "23505" {
				return User{}, &UsernameIsTakenError{}
			}
		}

		return User{}, errors.New("unknown error when trying to DELETE user by ID")
	}

	return user, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (User, error) {
	var user User
	err := r.pool.QueryRow(ctx, getUserByUsernameQuery, username).
		Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.Deleted,
		)

	if err != nil {
		slog.Error("[UserRepository-GetUserByUsername]", "Error", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, &UserDoesNotExistError{}
		}

		return User{}, errors.New("unknown error when trying to GET user by username")
	}

	return user, nil
}

func (r *UserRepository) DeleteUserByUsername(ctx context.Context, username string) (User, error) {
	var user User
	err := r.pool.QueryRow(ctx, deleteUserByUsernameQuery, username).
		Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.Deleted,
		)
	if err != nil {
		slog.Error("[UserRepository-DeleteUserByUsername]", "Error", err)

		// Need to check other errors if there is any
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, &UserDoesNotExistError{}
		}

		return User{}, errors.New("unknown error when trying to DELETE user by username")
	}

	return user, nil
}

func (r *UserRepository) UpdateUserByUsername(ctx context.Context, userId *string, newUsername *string, newEmail *string) (User, error) {
	var setClauses []string
	var args []interface{}
	paramIndex := 1

	if newUsername != nil {
		setClauses = append(setClauses, fmt.Sprintf("username = $%d", paramIndex))
		args = append(args, *newUsername)
		paramIndex++
	}
	if newEmail != nil {
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", paramIndex))
		args = append(args, *newEmail)
		paramIndex++
	}

	setClauses = append(setClauses, "updated_at = NOW()")
	args = append(args, *userId)

	query := fmt.Sprintf(
		updateUserByUsernameQuery,
		strings.Join(setClauses, ", "),
		paramIndex,
	)

	var user User
	err := r.pool.QueryRow(ctx, query, args...).
		Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.Deleted,
		)

	if err != nil {
		slog.Error("[UserRepository-UpdateUserByUsername]", "Error", err)
		// Need to check other errors if there is any
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, &UserDoesNotExistError{}
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// Duplicate key value violates unique constraint
			if pgErr.Code == "23505" {
				return User{}, &UsernameIsTakenError{}
			}
		}

		return User{}, errors.New("unknown error when trying to DELETE user by username")
	}

	return user, nil
}
