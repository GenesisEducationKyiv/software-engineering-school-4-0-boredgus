package db

import (
	"context"
	"fmt"
	"subscription-api/internal/services"
)

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

const createUserQ string = `
	insert into subs."users" (email)
	values ($1);
`

func (r *UserRepo) CreateUser(ctx context.Context, d DB, email string) error {
	_, err := d.DB().ExecContext(ctx, createUserQ, email)
	if d.IsError(err, UniqueViolation) {
		return fmt.Errorf("%w: user with such email already exists", services.UniqueViolationErr)
	}

	return err
}
