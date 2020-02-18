package tea

import (
	"fmt"
	"net/http"
	"time"
)

// // Headers map http headers
// type HTTPHeaders map[string]string

// Meta struct meta
type Meta struct {
	HTTP  HTTP
	Start Duration
	End   Duration
}

// HTTP struct for collecting HTTP status code and headers
type HTTP struct {
	Code    int
	Headers map[string]string
}

// Duration struct Duration
type Duration struct {
	Valid bool
	Time  time.Time
}

func headers(response *http.Response) map[string]string {
	tmp := make(map[string]string)
	for k, v := range response.Header {
		if 1 == len(v) {
			tmp[k] = fmt.Sprintf("%s", v[0])
		}
	}

	return tmp
}
