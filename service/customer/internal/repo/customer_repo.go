package repo

import (
	"strings"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/service"
	"github.com/google/uuid"
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
		Where(query any, args ...any) DB
	}

	customerRepo struct {
		db DB
	}
)

type DBError int

const (
	ErrUniqueViolation DBError = iota
)

func NewCustomerRepo(db DB) *customerRepo {
	return &customerRepo{
		db: db,
	}
}

func (r *customerRepo) CreateCustomer(email string) error {
	err := r.db.Create(Customer{
		ID:        uuid.NewString(),
		Email:     email,
		CreatedAt: time.Now().UTC(),
	})
	if r.isErr(err, ErrUniqueViolation) {
		return service.ErrAlreadyExists
	}

	return err
}

func (r *customerRepo) DeleteCustomer(email string) error {
	return r.db.Where("email = ?", email).Delete(&Customer{})
}

var errorMatch = map[DBError]string{
	ErrUniqueViolation: "23505",
}

func (r *customerRepo) isErr(err error, targetErr DBError) bool {
	if err == nil {
		return false
	}
	strToFind, ok := errorMatch[targetErr]
	if !ok {
		return false
	}

	return strings.Contains(err.Error(), strToFind)
}
