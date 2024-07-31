package repo

import (
	"strings"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/customer/internal/service"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	Customer struct {
		ID        string `gorm:"primaryKey"`
		Email     string `gorm:"unique"`
		CreatedAt time.Time
	}

	customerRepo struct {
		db *gorm.DB
	}
)

type DBError int

const (
	ErrUniqueViolation DBError = iota
)

func NewCustomerRepo(db *gorm.DB) *customerRepo {
	return &customerRepo{
		db: db,
	}
}

func (r *customerRepo) CreateCustomer(email string) error {
	err := r.db.Create(Customer{
		ID:        uuid.NewString(),
		Email:     email,
		CreatedAt: time.Now().UTC(),
	}).Error
	if r.isErr(err, ErrUniqueViolation) {
		return service.ErrAlreadyExists
	}

	return err
}

func (r *customerRepo) DeleteCustomer(email string) error {
	return r.db.Where("email = ?", email).Delete(&Customer{}).Error
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
