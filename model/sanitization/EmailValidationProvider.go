package sanitization

import "net/mail"

/*
EmailVaidationProvider is an interface for describing an email
validation service.
*/
type EmailValidationProvider interface {
	GetEmailComponents(email string) (*mail.Address, error)
	IsValidEmail(email string) bool
}
