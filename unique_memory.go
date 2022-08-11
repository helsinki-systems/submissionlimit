package submissionlimit

import (
	"fmt"
	"sync"
)

type uniqueMemoryStorage struct {
	store   map[string]map[string]struct{}
	storeMu sync.Mutex
}

var _ UniqueStorage = (*uniqueMemoryStorage)(nil)

func (ums *uniqueMemoryStorage) Store(k, v string) error {
	ums.storeMu.Lock()
	defer ums.storeMu.Unlock()

	if svs, ok := ums.store[k]; !ok {
		s := make(map[string]struct{})
		s[v] = struct{}{}
		ums.store[k] = s
	} else {
		if _, ok := svs[v]; ok {
			return fmt.Errorf("key %q not unique", k)
		}

		ums.store[k][v] = struct{}{}
	}

	return nil
}

func NewUniqueMemoryStorage() *uniqueMemoryStorage {
	return &uniqueMemoryStorage{
		store: make(map[string]map[string]struct{}),
	}
}
