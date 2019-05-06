package domain

import "context"

// NewCallbackPersisterMock creates a callback mock for tests
// nolint:unused
func NewCallbackPersisterMock(persistFunc func(g *Geolocation) error) Persister {
	return &storagePersisterMock{
		persistFunc: persistFunc,
	}
}

// nolint:unused
type storagePersisterMock struct {
	persistFunc func(g *Geolocation) error
}

func (m *storagePersisterMock) Persist(_ context.Context, g *Geolocation) error {
	return m.persistFunc(g)
}
