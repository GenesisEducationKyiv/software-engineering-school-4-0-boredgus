package repo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/entities"
)

type Dispatches map[string]entities.Dispatch

func (d Dispatches) ToReader() (io.Reader, error) {
	marshalled, err := json.Marshal(&d)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal dispatches: %w", err)
	}

	return bytes.NewReader(marshalled), nil
}

func SerializeDispatch(d entities.Dispatch) ([]byte, error) {
	return json.Marshal(d)
}

func DeserializeDispatches(data []byte) (map[string]entities.Dispatch, error) {
	var dsptches map[string]entities.Dispatch
	if err := json.Unmarshal(data, &dsptches); err != nil {
		return nil, fmt.Errorf("failed to deserialize dispatch: %w", err)
	}

	return dsptches, nil
}
