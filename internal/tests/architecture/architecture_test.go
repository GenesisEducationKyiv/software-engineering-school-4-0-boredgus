package architecture

import (
	"testing"

	"github.com/matthewmcnew/archtest"
)

func TestArchitectureDependencyViolationFor(t *testing.T) {
	externalInterfaces := []string{
		"go.uber.org/zap",
		"github.com/gin-gonic/gin",
		"github.com/lib/pq",
	}

	tests := []struct {
		name                         string
		basePackage                  string
		ignoredDependencies          []string
		notAllowedDependencies       []string
		notAllowedDirectDependencies []string
	}{
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
			notAllowedDirectDependencies: externalInterfaces,
		},
		{
			name:        "clients",
			basePackage: "subscription-api/internal/clients/...",
			notAllowedDependencies: []string{
				"subscription-api/cmd/...",
				"subscription-api/internal/controllers/...",
				"subscription-api/internal/services/dispatch/grpc",
				"subscription-api/internal/services/currency/grpc",
			},
			notAllowedDirectDependencies: externalInterfaces,
		},
		{
			name:        "services",
			basePackage: "subscription-api/internal/services",
			notAllowedDependencies: []string{
				"subscription-api/cmd/...",
				"subscription-api/internal/controllers/...",
				"subscription-api/internal/services/dispatch/grpc",
				"subscription-api/internal/services/currency/grpc",
			},
			notAllowedDirectDependencies: externalInterfaces,
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
			notAllowedDirectDependencies: externalInterfaces,
		},
		{
			name:        "dispatch service server",
			basePackage: "subscription-api/internal/services/dispatch/server",
			notAllowedDependencies: []string{
				"subscription-api/cmd/...",
				"subscription-api/internal/controllers/...",
			},
			notAllowedDirectDependencies: externalInterfaces,
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
			notAllowedDirectDependencies: externalInterfaces,
		},
		{
			name:        "dispatch currency server",
			basePackage: "subscription-api/internal/services/currency/server",
			notAllowedDependencies: []string{
				"subscription-api/cmd/...",
				"subscription-api/internal/controllers/...",
			},
			notAllowedDirectDependencies: externalInterfaces,
		},
		{
			name:        "controllers",
			basePackage: "subscription-api/internal/controllers",
			notAllowedDependencies: []string{
				"subscription-api/cmd/...",
				"subscription-api/internal/services/dispatch",
				"google.golang.org/grpc",
				"github.com/gin-gonic/gin",
			},
			notAllowedDirectDependencies: externalInterfaces,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockedT := new(testingT)

			packageTest := archtest.Package(mockedT, tt.basePackage).
				Ignoring(tt.ignoredDependencies...)

			packageTest.ShouldNotDependOn(tt.notAllowedDependencies...)
			packageTest.ShouldNotDependDirectlyOn(tt.notAllowedDirectDependencies...)

			mockedT.AssertNoError(t, mockedT)
		})
	}
}
