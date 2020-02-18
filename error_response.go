package tea

import "fmt"

// APIError struct APIError
type APIError struct {
	Errors []string `json:"errors"`
}

func (e APIError) Error() string {
	if len(e.Errors) > 0 {
		err := e.Errors[0]
		return fmt.Sprintf("twitter: %v", err)
	}
	return ""
}
