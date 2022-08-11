package submissionlimit

import (
	"errors"
	"fmt"
	"log"
	"net"
)

type Submission struct {
	IP string

	EMail string

	Unique map[string]string
}

func (l *Limiter) Limit(s Submission) error {
	// TODO: Ensure our own keys are unset
	uniq := func(k, v string) {
		if s.Unique == nil {
			s.Unique = make(map[string]string)
		}

		s.Unique[k] = v
	}

	if s.IP != "" {
		if l.ipLimiter == nil {
			log.Printf("warning: IP provided but IP limiter not set\n")
		} else {
			if l.ipLimiter.unique {
				uniq("ip", s.IP)
			}

			if err := l.isIPLimited(s.IP); err != nil {
				return fmt.Errorf("limited due to IP: %w", err)
			}
		}
	}

	if s.EMail != "" {
		if l.emailLimiter == nil {
			log.Printf("warning: email provided but email limiter not set\n")
		} else {
			if l.emailLimiter.unique {
				uniq("email", s.EMail)
			}

			if err := l.isEMailLimited(s.EMail); err != nil {
				return fmt.Errorf("limited due to email: %w", err)
			}
		}
	}

	if s.Unique != nil {
		if l.uniqueLimiter == nil {
			log.Printf("warning: unique provided but unique limiter not set\n")
		} else {
			if err := l.unique(s.Unique); err != nil {
				return fmt.Errorf("limited due to unique: %w", err)
			}
		}
	}

	return nil
}

func (l *Limiter) isIPLimited(ip string) error {
	if l.ipLimiter != nil {
		pip := net.ParseIP(ip)
		if pip == nil {
			return errors.New("not an IP")
		}

		if isIPv4 := pip.To4() != nil; isIPv4 && !l.ipLimiter.allowIPv4 {
			return errors.New("IPv4 not allowed")
		} else if !isIPv4 && !l.ipLimiter.allowIPv6 {
			return errors.New("IPv6 not allowed")
		}
	}

	if l.torLimiter != nil {
		if err := l.torLimiter.IsLimited(ip); err != nil {
			return err
		}
	}

	return nil
}

func (l *Limiter) isEMailLimited(email string) error {
	if l.emailLimiter != nil {
		return l.emailLimiter.IsLimited(email)
	}

	return nil
}

func (l *Limiter) unique(m map[string]string) error {
	if l.uniqueLimiter != nil {
		return l.uniqueLimiter.Unique(m)
	}

	return nil
}
