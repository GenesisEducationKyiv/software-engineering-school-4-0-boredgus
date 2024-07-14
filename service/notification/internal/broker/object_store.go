package broker

import (
	"context"
	"errors"
	"io"

	"github.com/nats-io/nats.go/jetstream"
)

var (
	ErrNotFound = errors.New("not found")

	mappedErrors = map[error]error{
		jetstream.ErrObjectNotFound: ErrNotFound,
	}
)

type ObjectStore struct {
	store jetstream.ObjectStore
}

func NewObjectStore(store jetstream.ObjectStore) *ObjectStore {
	return &ObjectStore{store: store}
}

func (s *ObjectStore) Get(ctx context.Context, name string) (io.ReadCloser, error) {
	r, err := s.store.Get(ctx, name)
	if err != nil {
		if mappedErrors[err] != nil {
			return nil, mappedErrors[err]
		}

		return nil, err
	}

	return r, err
}

func (s *ObjectStore) Put(ctx context.Context, name string, data io.Reader) error {
	_, err := s.store.Put(ctx, jetstream.ObjectMeta{Name: name}, data)

	return err
}
