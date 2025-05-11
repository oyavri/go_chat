package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	// go:embed sql/create_user.sql
	createUserQuery string
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) CreateUser(ctx context.Context, input CreateUserRequest) (User, error) {
	var user User

	err := r.pool.QueryRow(ctx, createUserQuery, input.Username, input.Email).Scan(&user)
	if err != nil {
		// Need to check other errors, it might be due to unique username
		return User{}, &UserAlreadyExistsError{}
	}

	return user, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, userId string) (User, error) {
	var user User
	query := `SELECT * FROM user
			  WHERE id = $1`

	err := r.pool.QueryRow(ctx, query, userId).Scan(&user)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, &UserDoesNotExistError{}
		}

		return User{}, errors.New("unknown error when trying to GET user by ID")
	}

	return user, nil
}

func (r *UserRepository) DeleteUserById(ctx context.Context, userId string) (User, error) {
	var user User
	query := `UPDATE user
			  SET deleted = true, deleted_at = NOW()
			  WHERE id = $1
			  RETURNING *`

	err := r.pool.QueryRow(ctx, query, user.Username, user.Email).Scan(&user)
	if err != nil {
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

	if len(setClauses) == 0 {
		return User{}, errors.New("no fields to update")
	}

	setClauses = append(setClauses, "updated_at = NOW()")
	args = append(args, *userId)

	query := fmt.Sprintf(
		`UPDATE user 
	     SET %s 
		 WHERE id = $%d 
		 RETURNING *`,
		strings.Join(setClauses, ", "),
		paramIndex,
	)

	var user User
	err := r.pool.QueryRow(ctx, query, args...).Scan(&user)
	if err != nil {
		// Need to check other errors if there is any
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, &UserDoesNotExistError{}
		}
		// There might be a conflict on username
		return User{}, errors.New("unknown error when trying to DELETE user by ID")
	}

	return user, nil
}
