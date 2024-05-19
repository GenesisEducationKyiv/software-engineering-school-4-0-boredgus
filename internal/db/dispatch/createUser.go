package ds

import (
	"context"
	"fmt"
	"subscription-api/internal/db"
	"subscription-api/internal/services"
)

const createUserQ string = `
	insert into subs."users" (email)
	values ($1);
`

func (s *currencyDispatchStore) CreateSubscriber(ctx context.Context, email string) error {
	_, err := s.DB().ExecContext(ctx, createUserQ, email)
	if s.IsError(err, db.UniqueViolation) {
		return fmt.Errorf("%w: user with such email already exists", services.UniqueViolationErr)
	}
	return err
}
