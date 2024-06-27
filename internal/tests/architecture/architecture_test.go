package architecture

import (
	"testing"

	"github.com/matthewmcnew/archtest"
)

func TestArchitectureDependencyViolationFor(t *testing.T) {
	tests := []struct {
		name                   string
		basePackage            string
		ignoredDependencies    []string
		notAllowedDependencies []string
	}{
		{
			name:        "services",
			basePackage: "subscription-api/internal/services",
			notAllowedDependencies: []string{
				"subscription-api/cmd/...",
				"subscription-api/internal/controllers/...",
				"subscription-api/internal/services/dispatch/grpc",
				"subscription-api/internal/services/currency/grpc",
			},
		},
		{
			name:        "dispatch service",
			basePackage: "subscription-api/internal/services/dispatch",
			notAllowedDependencies: []string{
				"subscription-api/cmd/...",
				"subscription-api/internal/controllers/...",
				"subscription-api/internal/services/dispatch/server",
				"subscription-api/internal/services/dispatch/server/grpc",
			},
		},
		{
			name:        "dispatch service server",
			basePackage: "subscription-api/internal/services/dispatch/server",
			notAllowedDependencies: []string{
				"subscription-api/cmd/...",
				"subscription-api/internal/controllers/...",
			},
		},
		{
			name:        "currency service",
			basePackage: "subscription-api/internal/services/currency",
			notAllowedDependencies: []string{
				"subscription-api/cmd/...",
				"subscription-api/internal/controllers/...",
				"subscription-api/internal/services/currency/server",
				"subscription-api/internal/services/currency/server/grpc",
			},
		},
		{
			name:        "dispatch currency server",
			basePackage: "subscription-api/internal/services/currency/server",
			notAllowedDependencies: []string{
				"subscription-api/cmd/...",
				"subscription-api/internal/controllers/...",
			},
		},
		{
			name:        "controllers",
			basePackage: "subscription-api/internal/controllers",
			notAllowedDependencies: []string{
				"subscription-api/cmd/...",
				"subscription-api/internal/services/dispatch",
			},
		},
		{
			name:        "entities",
			basePackage: "subscription-api/internal/entities",
			notAllowedDependencies: []string{
				"subscription-api/cmd/...",
				"subscription-api/internal/clients/...",
				"subscription-api/internal/controllers/...",
				"subscription-api/internal/mailing/...",
				"subscription-api/internal/db/...",
				"subscription-api/internal/mocks/...",
				"subscription-api/internal/services/...",
				"subscription-api/internal/sql/...",
				"subscription-api/internal/tests/...",
				"subscription-api/internal/tests/...",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockedT := new(testingT)

			archtest.Package(mockedT, tt.basePackage).
				Ignoring(tt.ignoredDependencies...).
				ShouldNotDependOn(tt.notAllowedDependencies...)

			mockedT.AssertNoError(t, mockedT)
		})
	}
}
