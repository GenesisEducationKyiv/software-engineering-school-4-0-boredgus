package stubs

import (
	"context"
	"testing"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/service"
	"github.com/stretchr/testify/assert"
)

type NotificationServiceMock struct {
	calls []service.NotificationType
}

func NewNotificationServiceMock() *NotificationServiceMock {
	return &NotificationServiceMock{}
}

func (m *NotificationServiceMock) AssertNotificationTypeOfLastCall(t *testing.T, nType service.NotificationType) {
	t.Helper()

	assert.Equal(t, nType, m.calls[len(m.calls)-1])
}

func (m *NotificationServiceMock) SendSubscriptionDetails(ctx context.Context, notification service.Notification) error {
	m.calls = append(m.calls, notification.Type)

	return nil
}
func (m *NotificationServiceMock) SendExchangeRates(ctx context.Context, notification service.Notification) error {
	m.calls = append(m.calls, notification.Type)

	return nil
}
