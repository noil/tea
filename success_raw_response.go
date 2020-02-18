package tea

import "net/http"

// SuccessRaw map SuccessRaw
type SuccessRaw map[string]interface{}

// ResponseRaw struct ResponseRaw
type ResponseRaw struct {
	Meta Meta
	Data SuccessRaw
}

func newResponseRaw() *ResponseRaw {
	return &ResponseRaw{
		Data: make(SuccessRaw),
		Meta: Meta{HTTP: HTTP{Headers: make(map[string]string)}},
	}
}

func (result *ResponseRaw) meta(response *http.Response) {
	if nil != response {
		result.Meta.HTTP.Code = response.StatusCode
		result.Meta.HTTP.Headers = headers(response)
	}
}
