package broker

import (
	"context"
	"io"

	"github.com/nats-io/nats.go/jetstream"
)

type objectStore struct {
	store jetstream.ObjectStore
}

func NewObjectStore(store jetstream.ObjectStore) *objectStore {
	return &objectStore{store: store}
}

func (s *objectStore) Get(ctx context.Context, name string) (io.ReadCloser, error) {
	return s.store.Get(ctx, name)
}

func (s *objectStore) Put(ctx context.Context, name string, data io.Reader) error {
	_, err := s.store.Put(ctx, jetstream.ObjectMeta{Name: name}, data)

	return err
}
