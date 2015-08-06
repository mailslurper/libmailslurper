package sanitization

/*
XSSServiceProvider is an interface for providing cross-site scripting
and sanitization services.
*/
type XSSServiceProvider interface {
	SanitizeString(input string) string
}
