package services

import "errors"

var (
	UniqueViolationErr    = errors.New("unique violation")
	InvalidArgumentErr    = errors.New("invalid argument")
	NotFoundErr           = errors.New("not found")
	FailedPreconditionErr = errors.New("failed precondition")
)

const USD_UAH_DISPATCH_ID = "f669a90d-d4aa-4285-bbce-6b14c6ff9065"
