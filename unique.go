package submissionlimit

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
