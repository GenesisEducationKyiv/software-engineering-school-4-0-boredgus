package notifier_test

import (
	"testing"

	notifier_mock "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/mocks/notifier"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/service"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/service/notifier"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/service/notifier/emails"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/service/notifier/emails/templates"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_EmailNotifier_Notify(t *testing.T) {
	type args struct {
		notification service.Notification
	}

	mailmanMock := notifier_mock.NewMailman(t)

	tests := []struct {
		name        string
		args        args
		setup       func(args) func()
		expectedErr error
	}{
		{
			name: "skipped: there is no provided emails",
			args: args{},
			setup: func(a args) func() {
				return func() {}
			},
			expectedErr: nil,
		},
		{
			name: "failed: unsupported notification type",
			args: args{notification: service.Notification{
				Type: -2,
				Data: service.NotificationData{Emails: []string{"test@gmail.com"}},
			}},
			setup: func(a args) func() {
				return func() {}
			},
			expectedErr: templates.UnknownNotificationTypeErr,
		},
		{
			name: "failed: payload does not suite email template",
			args: args{notification: service.Notification{
				Type: service.SubscriptionCreated,
				Data: service.NotificationData{
					Emails:  []string{"test@gmail.com"},
					Payload: struct{ s int }{s: 12},
				},
			}},
			setup: func(a args) func() {
				return func() {}
			},
			expectedErr: emails.HTMLTemplateErr,
		},
		{
			name: "failed:got an error while sending email",
			args: args{notification: service.Notification{
				Type: service.SubscriptionCreated,
				Data: service.NotificationData{
					Emails: []string{"test@gmail.com"},
				},
			}},
			setup: func(a args) func() {
				sendCall := mailmanMock.EXPECT().
					Send(mock.Anything).Return(assert.AnError)

				return func() {
					sendCall.Unset()
				}
			},
			expectedErr: assert.AnError,
		},
		{
			name: "success: sent an email",
			args: args{notification: service.Notification{
				Type: service.SubscriptionCreated,
				Data: service.NotificationData{
					Emails: []string{"test@gmail.com"},
				},
			}},
			setup: func(a args) func() {
				sendCall := mailmanMock.EXPECT().Send(mock.Anything).Return(nil)

				return func() {
					sendCall.Unset()
				}
			},
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setup(tt.args)
			defer cleanup()

			n := notifier.NewEmailNotifier(mailmanMock)
			actualErr := n.Notify(tt.args.notification)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, actualErr, tt.expectedErr)

				return
			}
			assert.Nil(t, actualErr)
		})
	}
}
