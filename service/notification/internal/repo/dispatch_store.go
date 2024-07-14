package repo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/broker"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/entities"
)

const DispathesObjectName string = "scheduled_dispatches"

type (
	dispatchStore struct {
		store      *broker.ObjectStore
		mu         *sync.Mutex
		dispatches map[string]entities.Dispatch
	}
)

func NewDispatchRepo(store *broker.ObjectStore) *dispatchStore {
	return &dispatchStore{
		store:      store,
		mu:         &sync.Mutex{},
		dispatches: make(map[string]entities.Dispatch),
	}
}

func (s *dispatchStore) AddSubscription(ctx context.Context, sub *entities.Subscription) error {
	s.mu.Lock()
	dispatch, ok := s.dispatches[sub.DispatchID]
	if !ok {
		s.dispatches[sub.DispatchID] = *sub.ToDispatch()
	} else {
		dispatch.Emails = append(s.dispatches[sub.DispatchID].Emails, sub.Email)
		s.dispatches[sub.DispatchID] = dispatch
	}
	s.mu.Unlock()

	reader, err := Dispatches(s.dispatches).ToReader()
	if err != nil {
		return fmt.Errorf("failed to write dispatches into io.Reader: %w", err)
	}

	err = s.store.Put(ctx, DispathesObjectName, reader)
	if err != nil {
		return fmt.Errorf("failed to updated dispatches: %w", err)
	}

	return nil
}

func (s *dispatchStore) GetAll(ctx context.Context) (map[string]entities.Dispatch, error) {
	reader, err := s.store.Get(ctx, DispathesObjectName)
	if errors.Is(err, broker.ErrNotFound) {
		reader, err := Dispatches(s.dispatches).ToReader()
		if err != nil {
			return nil, fmt.Errorf("failed to write dispatches into io.Reader: %w", err)
		}

		if err = s.store.Put(ctx, DispathesObjectName, reader); err != nil {
			return nil, fmt.Errorf("failed to set empty sset of dispatches: %w", err)
		}

		return s.dispatches, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get dispatches data: %w", err)
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(reader); err != nil {
		return nil, fmt.Errorf("failed to read dispatches: %w", err)
	}

	var dispatches map[string]entities.Dispatch
	if err := json.Unmarshal(buf.Bytes(), &dispatches); err != nil {
		return nil, fmt.Errorf("failed to deserialize dispatches: %w", err)
	}

	return dispatches, nil
}
