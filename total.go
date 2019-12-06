package tea

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (tea *TwitterEngagementAPI) Total(tweetIds []string, types []EngagementType, groupBy []GroupBy) {
	params, err := tea.Params(tweetIds, types, groupBy)
	if err != nil {
		fmt.Println(err)
	}
	client := tea.httpClient()
	request, err := http.NewRequest("POST", totalUrl, bytes.NewBuffer(params))
	request.Header.Add("Accept-Encoding", "gzip")
	request.Header.Add("Authorization", tea.OAuthHeader(totalUrl))
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)
	defer func() {
		if nil != response {
			response.Body.Close()
		}
	}()
	if err != nil {
		panic(err)
	}
	if http.StatusBadRequest == response.StatusCode {
		panic("bad http code")
	}
	reader, err := gzip.NewReader(response.Body)
	if nil != err {
		panic(err)
	}
	data, err := ioutil.ReadAll(reader)
	if nil != err {
		panic(err)
	}
	fmt.Printf("%s", data)
	// err = json.NewDecoder(response.Body).Decode(token)
	// if nil != err {
	// 	return token, err
	// }

	// fmt.Printf("%s\r\n", params)
	// fmt.Printf("%s\r\n", sign)
}

func (tea *TwitterEngagementAPI) Params(tweetIds []string, types []EngagementType, groupBy []GroupBy) ([]byte, error) {
	return requestParams(tweetIds, types, groupBy)
}
