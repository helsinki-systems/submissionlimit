package submissionlimit_test

import (
	"testing"

	sl "github.com/helsinki-systems/submissionlimit"
)

func TestIPLimiter(t *testing.T) {
	t.Parallel()

	l := sl.New(
		sl.WithIP(sl.IPConfig{
			AllowIPv4: true,
		}),
	)

	if err := l.Limit(sl.Submission{
		IP: "1.2.3.4",
	}); err != nil {
		t.Errorf("submission with IP 1.2.3.4 should not be limited: %v", err)
	}
	if err := l.Limit(sl.Submission{
		IP: "1.2.3.4",
	}); err != nil {
		t.Errorf("submission with IP 1.2.3.4 should not be limited: %v", err)
	}

	if err := l.Limit(sl.Submission{
		IP: "fe80::",
	}); err == nil {
		t.Error("submission with IP fe80:: should be limited")
	}
}

func TestEMailLimiter(t *testing.T) {
	t.Parallel()

	l := sl.New(
		sl.WithEMail(sl.EMailConfig{}),
	)

	if err := l.Limit(sl.Submission{
		EMail: "a",
	}); err == nil {
		t.Error("submission with invalid email should be limited")
	}

	if err := l.Limit(sl.Submission{
		EMail: "hello@example.com",
	}); err != nil {
		t.Errorf("submission with valid email should not be limited: %v", err)
	}

	if err := l.Limit(sl.Submission{
		EMail: "hello@throwawayemail.com",
	}); err == nil {
		t.Error("submission with disposable email should be limited")
	}
}

func TestUniqueLimiter(t *testing.T) {
	t.Parallel()

	l := sl.New(
		sl.WithIP(sl.IPConfig{
			AllowIPv4: true,

			Unique: true,
		}),

		sl.WithUnique(sl.UniqueConfig{}),
	)

	if err := l.Limit(sl.Submission{
		IP: "1.2.3.4",
	}); err != nil {
		t.Errorf("submission with IP 1.2.3.4 should not be limited: %v", err)
	}
	if err := l.Limit(sl.Submission{
		IP: "1.2.3.4",
	}); err == nil {
		t.Error("submission with IP 1.2.3.4 should now be limited due to unique")
	}

	if err := l.Limit(sl.Submission{
		Unique: map[string]string{
			"hello": "world",
		},
	}); err != nil {
		t.Errorf("submission with never seen unique key should not be limited: %v", err)
	}
	if err := l.Limit(sl.Submission{
		Unique: map[string]string{
			"hello": "world2",
		},
	}); err != nil {
		t.Errorf("submission with never seen unique key should not be limited: %v", err)
	}
	if err := l.Limit(sl.Submission{
		Unique: map[string]string{
			"hello": "world",
		},
	}); err == nil {
		t.Error("submission with repeated unique key should be limited")
	}
}

type nopStore struct {
	store map[string]map[string]struct{}
}

var _ sl.UniqueStorage = (*nopStore)(nil)

func (ns *nopStore) Store(k, v string) error {
	return nil
}

func TestUniqueLimiterWithNOPStore(t *testing.T) {
	t.Parallel()

	l := sl.New(
		sl.WithIP(sl.IPConfig{
			AllowIPv4: true,

			Unique: true,
		}),

		sl.WithUnique(sl.UniqueConfig{
			Storage: &nopStore{
				store: make(map[string]map[string]struct{}),
			},
		}),
	)

	if err := l.Limit(sl.Submission{
		IP: "1.2.3.4",
	}); err != nil {
		t.Errorf("submission with IP 1.2.3.4 should not be limited: %v", err)
	}
	if err := l.Limit(sl.Submission{
		IP: "1.2.3.4",
	}); err != nil {
		t.Errorf("submission with IP 1.2.3.4 should not be limited: %v", err)
	}
}
