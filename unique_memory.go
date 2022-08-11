package submissionlimit

import (
	"fmt"
	"sync"
)

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
