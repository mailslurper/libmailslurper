package sanitization

import (
	"github.com/microcosm-cc/bluemonday"
)

/*
XSSService implements the XSSServiceProvider interface and offers functions to
help address cross-site script and sanitization concerns.
*/
type XSSService struct {
	sanitizer *bluemonday.Policy
}

/*
NewXSSService creates a new cross-site scripting service.
*/
func NewXSSService() *XSSService {
	return &XSSService{
		sanitizer: bluemonday.UGCPolicy(),
	}
}

/*
SanitizeString attempts to sanitize a string by removing potentially dangerous
HTML/JS markup.
*/
func (service *XSSService) SanitizeString(input string) string {
	return service.sanitizer.Sanitize(input)
}
