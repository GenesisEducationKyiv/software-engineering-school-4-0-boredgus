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
	}
	customerService struct {
		repo CustomerRepo
	}
)

var AlreadyExistsErr = errors.New("already exists")

func NewCustomerService(repo CustomerRepo) *customerService {
	return &customerService{
		repo: repo,
	}
}

func (s *customerService) CreateCustomer(ctx context.Context, email string) error {
	return nil
}
