package submissionlimit

import (
	"errors"
	"fmt"

	mailchecker "github.com/FGRibreau/mailchecker/v4/platform/go"
	"github.com/badoux/checkmail"
)

var (
	ErrEMailInvalidOrOnDenyList = errors.New("email invalid or domain on deny list")
)

type EMailConfig struct {
	AdditionalDenyDomains []string

	DoValidateDomain        bool
	DoValidateDomainAndUser bool

	SMTPHostname string
	SMTPEMail    string

	Unique bool
}

type emailLimiter struct {
	doValidateDomain        bool
	doValidateDomainAndUser bool

	smtpHostname string
	smtpEMail    string

	unique bool
}

func (el *emailLimiter) IsLimited(email string) error {
	if !mailchecker.IsValid(email) {
		return ErrEMailInvalidOrOnDenyList
	}

	if el.doValidateDomainAndUser {
		if err := checkmail.ValidateHostAndUser(
			el.smtpHostname,
			el.smtpEMail,
			email,
		); err != nil {
			return fmt.Errorf("failed to validate domain and user: %w", err)
		}
	} else if el.doValidateDomain {
		if err := checkmail.ValidateHost(email); err != nil {
			return fmt.Errorf("failed to validate domain: %w", err)
		}
	}

	return nil
}

func WithEMail(ec EMailConfig) option {
	if len(ec.AdditionalDenyDomains) != 0 {
		_ = mailchecker.AddCustomDomains(ec.AdditionalDenyDomains)
	}

	if ec.DoValidateDomainAndUser && (ec.SMTPHostname == "" || ec.SMTPEMail == "") {
		panic(errors.New("email: DoValidateDomainAndUser set but no SMTP details provided"))
	}

	el := &emailLimiter{
		doValidateDomain:        ec.DoValidateDomain,
		doValidateDomainAndUser: ec.DoValidateDomainAndUser,

		smtpHostname: ec.SMTPHostname,
		smtpEMail:    ec.SMTPEMail,

		unique: ec.Unique,
	}

	return func(l *Limiter) {
		l.emailLimiter = el
	}
}
