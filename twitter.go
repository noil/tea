package tea

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

const (
	defaultGrouping     = "default_grouping"
	unsupportedTweetIds = "unsupported_for_impressions_engagements_tweet_ids"
	unavailableTweetIds = "unavailable_tweet_ids"
	start               = "start"
	end                 = "end"
)

// TwitterEngagementAPI provides access to Twitter Engagement API.
type TwitterEngagementAPI struct {
	ctx        context.Context
	httpClient *http.Client
	token      Token
}

var totalURL = "https://data-api.twitter.com/insights/engagement/totals"
var eng28hrURL = "https://data-api.twitter.com/insights/engagement/28hr"
var historicalURL = "https://data-api.twitter.com/insights/engagement/historical"

// Token struct token
type Token struct {
	Access            string
	AccessSecret      string
	ConsumerKey       string
	ConsumerKeySecret string
}

// New creates an instance of TwitterEngagementAPI
func New(token Token, httpClient *http.Client) *TwitterEngagementAPI {
	return NewWithContext(context.Background(), token, httpClient)
}

// NewWithContext creates an instance of TwitterEngagementAPI with context
func NewWithContext(ctx context.Context, token Token, httpClient *http.Client) *TwitterEngagementAPI {
	if nil == httpClient {
		httpClient = http.DefaultClient
	}
	return &TwitterEngagementAPI{ctx: ctx, token: token, httpClient: httpClient}
}

// Client provide to replace default http.Client
func (tea *TwitterEngagementAPI) Client(httpClient *http.Client) *TwitterEngagementAPI {
	if nil != httpClient {
		tea.httpClient = httpClient
	}

	return tea
}

// AccessToken update access and access secret tokens
func (tea *TwitterEngagementAPI) AccessToken(access, accessSecret string) *TwitterEngagementAPI {
	tea.token.Access = access
	tea.token.AccessSecret = accessSecret

	return tea
}

func (tea *TwitterEngagementAPI) do(url string, data interface{}) (*http.Response, error) {
	params, err := json.Marshal(data)
	request, err := http.NewRequestWithContext(tea.ctx, "POST", url, bytes.NewBuffer(params))
	request.Header.Add("Accept-Encoding", "gzip")
	request.Header.Add("Authorization", tea.oauthHeader(url))
	request.Header.Add("Content-Type", "application/json")
	response, err := tea.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
