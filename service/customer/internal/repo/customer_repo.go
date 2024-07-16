package repo

import (
	"context"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/service"
)

type (
	Customer struct {
		ID        string `gorm:"primaryKey"`
		Email     string `gorm:"unique"`
		CreatedAt time.Time
	}

	DB interface {
		Create(value interface{}) error
		Delete(value interface{}) error
	}

	customerRepo struct {
		db    DB
		isErr func(error, DBError) bool
	}
)

type DBError int

const UniqueViolationErr DBError = iota

func NewCustomerRepo(db DB, isErr func(error, DBError) bool) *customerRepo {
	return &customerRepo{
		db:    db,
		isErr: isErr,
	}
}

func (r *customerRepo) CreateCustomer(ctx context.Context, email string) error {
	err := r.db.Create(Customer{})
	if r.isErr(err, UniqueViolationErr) {
		return service.AlreadyExistsErr
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *customerRepo) DeleteCustomer(ctx context.Context, email string) error {
	return nil
}
