package sanitization

import "net/mail"

/*
EmailValidationService realizes the EmailValidationProvider
interface by offering functions for working with email validation
and manipulation.
*/
type EmailValidationService struct {
}

func (service *EmailValidationService) GetEmailComponents(email string) (*mail.Address, error) {
	return mail.ParseAddress(email)
}

func (service *EmailValidationService) IsValidEmail(email string) bool {
	_, err := service.GetEmailComponents(email)
	return err == nil
}

func NewEmailValidationService() *EmailValidationService {
	return &EmailValidationService{}
}
