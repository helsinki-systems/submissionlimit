package submissionlimit

type IPConfig struct {
	AllowIPv4 bool
	AllowIPv6 bool

	Unique bool
}

type ipLimiter struct {
	allowIPv4 bool
	allowIPv6 bool

	unique bool
}

func WithIP(ipc IPConfig) option {
	ipl := &ipLimiter{
		allowIPv4: ipc.AllowIPv4,
		allowIPv6: ipc.AllowIPv6,

		unique: ipc.Unique,
	}

	return func(l *Limiter) {
		l.ipLimiter = ipl
	}
}
