package tea

import (
	"fmt"
	"net/http"
)

type Headers map[string]string

type HTTPMeta struct {
	Code    int
	Headers Headers
}

func headers(response *http.Response) Headers {
	tmp := make(map[string]string)
	for k, v := range response.Header {
		if 1 == len(v) {
			tmp[k] = fmt.Sprintf("%s", v[0])
		}
	}

	return tmp
}
