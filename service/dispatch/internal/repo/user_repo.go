package repo

import (
	"context"
	"fmt"

	service "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service"
)

type UserRepo struct {
	db DB
}

func NewUserRepo(db DB) *UserRepo {
	return &UserRepo{db: db}
}

const createUserQ string = `
	insert into subs."users" (email)
	values ($1);
`

func (r *UserRepo) CreateUser(ctx context.Context, email string) error {
	_, err := r.db.ExecContext(ctx, createUserQ, email)
	if r.db.IsError(err, UniqueViolation) {
		return fmt.Errorf("%w: user with such email already exists", service.ErrUniqueViolation)
	}

	return err
}
