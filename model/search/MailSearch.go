package search

/*
MailSearch is a set of criteria used to filter a mail collection
*/
type MailSearch struct {
	Message string
	Start   string
	End     string
	From    string
	To      string
}
