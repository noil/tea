package tea

import "fmt"

// EngAPIError represents a Twitter API Error response
// https://developer.twitter.com/en/docs/metrics/get-tweet-engagement/api-reference/post-insights-engagement#ErrorMessages
type EngAPIError struct {
	Errors []string `json:"errors"`
}

func (e EngAPIError) Error() string {
	if len(e.Errors) > 0 {
		err := e.Errors[0]
		return fmt.Sprintf("twitter: %v", err)
	}
	return ""
}

// APIError represents a Twitter API Error response
// https://dev.twitter.com/overview/api/response-codes
type APIError struct {
	Errors []ErrorDetail `json:"errors"`
}

// ErrorDetail represents an individual item in an APIError.
type ErrorDetail struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e APIError) Error() string {
	if len(e.Errors) > 0 {
		err := e.Errors[0]
		return fmt.Sprintf("twitter: %d %v", err.Code, err.Message)
	}
	return ""
}
