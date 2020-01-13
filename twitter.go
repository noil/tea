package tea

import "net/http"

type TwitterEngagementAPI struct {
	Token      Token
	httpClient *http.Client
}

var totalUrl = "https://data-api.twitter.com/insights/engagement/totals"
var eng28hrUrl = "https://data-api.twitter.com/insights/engagement/28hr"
var historicalUrl = "https://data-api.twitter.com/insights/engagement/historical"

type Token struct {
	Access            string
	AccessSecret      string
	ConsumerKey       string
	ConsumerKeySecret string
}

func New(token Token) *TwitterEngagementAPI {
	return &TwitterEngagementAPI{Token: token, httpClient: defaultHttpClient()}
}

func (tea *TwitterEngagementAPI) Client(httpClient *http.Client) *TwitterEngagementAPI {
	if nil != httpClient {
		tea.httpClient = httpClient
	}

	return tea
}
