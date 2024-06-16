package mailing

import (
	"subscription-api/internal/mailing/emails/test"
	config_mocks "subscription-api/internal/mocks/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_htmlTemplateParser_Parse(t *testing.T) {
	type args struct {
		templateName string
		data         any
	}

	loggerMock := config_mocks.NewLogger(t)
	setup := func(a *args) func() {
		logCall := loggerMock.EXPECT().
			Errorf(mock.Anything, mock.Anything, mock.Anything)

		return func() {
			logCall.Unset()
		}
	}

	tests := []struct {
		name    string
		args    *args
		want    []byte
		wantErr bool
	}{
		{
			name: "failed to execute html template",
			args: &args{
				templateName: "test/test",
				data:         "",
			},
			wantErr: true,
		},
		{
			name: "failed to execute html template",
			args: &args{
				templateName: "test/test",
				data:         test.TestData{FirstField: "first", SecondField: 22},
			},
			want:    []byte("<div><span>first</span><span>22</span></div>"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setup(tt.args)
			defer cleanup()
			p := htmlTemplateParser{
				l: loggerMock,
			}
			got, err := p.Parse(tt.args.templateName, tt.args.data)

			assert.Equal(t, tt.want, got)
			if !tt.wantErr {
				assert.Nil(t, err, tt.wantErr)
			}
		})
	}
}
