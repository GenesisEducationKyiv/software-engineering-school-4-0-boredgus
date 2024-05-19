package services

import "errors"

var (
	UniqueViolationErr = errors.New("unique violation")
	InvalidArgumentErr = errors.New("invalid argument")
	NotFoundErr        = errors.New("not found")
)
