package broker

import (
	"context"
	"io"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/repo"
	"github.com/nats-io/nats.go/jetstream"
)

var mappedErrors = map[error]error{
	jetstream.ErrObjectNotFound: repo.ErrNotFound,
}

type objectStore struct {
	store jetstream.ObjectStore
}

func NewObjectStore(store jetstream.ObjectStore) *objectStore {
	return &objectStore{store: store}
}

func (s *objectStore) Get(ctx context.Context, name string) (io.ReadCloser, error) {
	r, err := s.store.Get(ctx, name)
	if err != nil {
		if mappedErrors[err] != nil {
			return nil, mappedErrors[err]
		}

		return nil, err
	}

	return r, err
}

func (s *objectStore) Put(ctx context.Context, name string, data io.Reader) error {
	_, err := s.store.Put(ctx, jetstream.ObjectMeta{Name: name}, data)

	return err
}
