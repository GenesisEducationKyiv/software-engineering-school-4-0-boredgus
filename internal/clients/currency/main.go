package currency_client

import (
	"encoding/json"
	"errors"
	"subscription-api/config"
)

var (
	UnsupportedCurrencyErr   = errors.New("unsupported currency")
	ServiceIsUnaccessibleErr = errors.New("service is unaccessible")
)

type (
	responseParams struct {
		Issuer     string         `json:"issuer,omitempty"`
		StatusCode int            `json:"status_code,omitempty"`
		Error      string         `json:"error,omitempty"`
		Data       map[string]any `json:"data,omitempty"`
	}

	responseLogger func(status int, err error, data map[string]any)
)

func (p responseParams) String() string {
	marshalledData, err := json.Marshal(p)
	if err != nil {
		return ""
	}

	return string(marshalledData)
}

func buildResponseLogger(logger config.Logger, issuer string) responseLogger {
	return func(status int, err error, data map[string]any) {
		params := responseParams{
			Issuer:     issuer,
			StatusCode: status,
			Data:       data,
		}

		if err != nil {
			params.Error = err.Error()
			logger.Error(params)

			return
		}
		logger.Info(params)
	}
}
