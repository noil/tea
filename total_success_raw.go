package tea

import "net/http"

type APISuccessRaw map[string]interface{}
type TotalAPISuccessRaw struct {
	Meta HTTPMeta
	Data APISuccessRaw
}

func newTotalAPISuccessRaw() *TotalAPISuccessRaw {
	return &TotalAPISuccessRaw{
		Data: make(APISuccessRaw),
		Meta: HTTPMeta{Headers: make(Headers)},
	}
}

func (result *TotalAPISuccessRaw) meta(response *http.Response) {
	if nil != response {
		result.Meta.Code = response.StatusCode
		result.Meta.Headers = headers(response)
	}
}
