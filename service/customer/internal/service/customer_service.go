package service

import (
	"context"
	"errors"
)

type (
	Customer struct {
		ID    string
		Email string
	}
	CustomerRepo interface {
		CreateCustomer(email string) error
		DeleteCustomer(email string) error
	}
	customerService struct {
		repo CustomerRepo
	}
)

var (
	ErrAlreadyExists = errors.New("already exists")
)

func NewCustomerService(repo CustomerRepo) *customerService {
	return &customerService{
		repo: repo,
	}
}

func (s *customerService) CreateCustomer(ctx context.Context, email string) error {
	return s.repo.CreateCustomer(email)
}

func (s *customerService) DeleteCustomer(ctx context.Context, email string) error {
	return s.repo.DeleteCustomer(email)
}
