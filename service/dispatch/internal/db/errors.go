package db

type Error int

const (
	UniqueViolation Error = iota + 1
	InvalidTextRepresentation
)

type ErrorCheckFunc func(error, Error) bool
