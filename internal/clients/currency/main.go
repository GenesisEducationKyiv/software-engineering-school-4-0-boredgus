package currency_client

import (
	"encoding/json"
	"errors"
)

var (
	UnsupportedCurrencyErr   = errors.New("unsupported currency")
	ServiceIsUnaccessibleErr = errors.New("service is unaccessible")
)

type responseParams struct {
	Issuer             string         `json:"issuer,omitempty"`
	ResponseStatusCode int            `json:"response_status_code,omitempty"`
	Error              string         `json:"error,omitempty"`
	Data               map[string]any `json:"data,omitempty"`
}

func (p responseParams) String() string {
	marshalledData, err := json.Marshal(p)
	if err != nil {
		return ""
	}

	return string(marshalledData)
}
