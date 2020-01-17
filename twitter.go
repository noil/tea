package tea

import (
	"context"
	"net/http"
)

//A TwitterEngagementAPI is a struct that provide access to Twitter Engagement API.
type TwitterEngagementAPI struct {
	ctx        context.Context
	httpClient *http.Client
	Token      Token
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

func New(token Token, httpClient *http.Client) *TwitterEngagementAPI {
	return NewWithContext(context.Background(), token, httpClient)
}

func NewWithContext(ctx context.Context, token Token, httpClient *http.Client) *TwitterEngagementAPI {
	if nil == httpClient {
		httpClient = http.DefaultClient
	}
	return &TwitterEngagementAPI{ctx: ctx, Token: token, httpClient: httpClient}
}

func (tea *TwitterEngagementAPI) Client(httpClient *http.Client) *TwitterEngagementAPI {
	if nil != httpClient {
		tea.httpClient = httpClient
	}

	return tea
}
