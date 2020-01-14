package tea

import "net/http"

type Success map[string]interface{}
type TotalResult struct {
	Data Success
	Meta HTTPMeta
}

func newTotalResult() *TotalResult {
	return &TotalResult{
		Data: make(Success),
		Meta: HTTPMeta{Headers: make(Headers)},
	}
}

func (result *TotalResult) meta(response *http.Response) {
	if nil != response {
		result.Meta.Code = response.StatusCode
		result.Meta.Headers = headers(response)
	}
}
