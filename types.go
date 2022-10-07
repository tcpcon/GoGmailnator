package GoGmailnator

import "fmt"

type Session struct {
	XsrfToken         string
	GmailnatorSession string
	Proxy             *string
}

type Email struct {
	From    string
	Subject string
	Time    string
}

type RetrieveEmailsJsonResponse struct {
	MessageData []Email
}

type GenerateEmailJsonResponse struct {
	Email []string
}

type RequestErr struct {
	StatusCode int
	Err        error
}

func (r *RequestErr) Error() string {
	return fmt.Sprintf("Status Code %d: Err %v", r.StatusCode, r.Err)
}
