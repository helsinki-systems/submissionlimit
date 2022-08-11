package submissionlimit_test

import (
	"testing"

	sl "github.com/helsinki-systems/submissionlimit"
)

func TestLimiter(t *testing.T) {
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
		IP: "fe80::",
	}); err == nil {
		t.Error("submission with IP fe80:: should be limited")
	}
}
