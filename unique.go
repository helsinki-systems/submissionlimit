package submissionlimit

import (
	"fmt"
	"sync"
)

type UniqueConfig struct {
	Storage uniqueStorage
}

type uniqueLimiter struct {
	storage uniqueStorage
}

type uniqueStorage interface {
	Store(k, v string) error
}

func (ul *uniqueLimiter) Unique(m map[string]string) error {
	for k, v := range m {
		if err := ul.storage.Store(k, v); err != nil {
			return err
		}
	}

	return nil
}

func WithUnique(uc UniqueConfig) option {
	if uc.Storage == nil {
		uc.Storage = NewUniqueMemoryStorage()
	}

	ul := &uniqueLimiter{
		storage: uc.Storage,
	}

	return func(l *Limiter) {
		l.uniqueLimiter = ul
	}
}

type UniqueMemoryStorage struct {
	store   map[string][]string
	storeMu sync.Mutex
}

var _ uniqueStorage = (*UniqueMemoryStorage)(nil)

func (ums *UniqueMemoryStorage) Store(k, v string) error {
	ums.storeMu.Lock()
	defer ums.storeMu.Unlock()

	if svs, ok := ums.store[k]; !ok {
		ums.store[k] = []string{v}
	} else {
		for _, sv := range svs {
			if sv == v {
				return fmt.Errorf("key %q not unique", k)
			}
		}

		ums.store[k] = append(ums.store[k], v)
	}

	return nil
}

func NewUniqueMemoryStorage() *UniqueMemoryStorage {
	return &UniqueMemoryStorage{
		store: make(map[string][]string),
	}
}
