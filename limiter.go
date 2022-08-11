package submissionlimit

type Limiter struct {
	ipLimiter     *ipLimiter
	emailLimiter  *emailLimiter
	torLimiter    *torLimiter
	uniqueLimiter *uniqueLimiter
}

func New(opts ...option) *Limiter {
	l := &Limiter{}

	for _, opt := range opts {
		opt(l)
	}

	return l
}
