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
	appOnly    bool
}

var (
	totalURL       = "https://data-api.twitter.com/insights/engagement/totals"
	eng28hrURL     = "https://data-api.twitter.com/insights/engagement/28hr"
	historicalURL  = "https://data-api.twitter.com/insights/engagement/historical"
	oauth2TokenURL = "https://api.twitter.com/oauth2/token"
)

// Token struct Token
type Token struct {
	Access            string
	AccessSecret      string
	ConsumerKey       string
	ConsumerKeySecret string
}

// AppOnlyToken struct AppOnlyToken
type AppOnlyToken struct {
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

// NewAppOnly creates an instance of TwitterEngagementAPI App Only
func NewAppOnly(token AppOnlyToken, httpClient *http.Client) *TwitterEngagementAPI {
	return NewAppOnlyWithContext(context.Background(), token, httpClient)
}

// NewAppOnlyWithContext creates an instance of TwitterEngagementAPI App Only with context
func NewAppOnlyWithContext(ctx context.Context, token AppOnlyToken, httpClient *http.Client) *TwitterEngagementAPI {
	if nil == httpClient {
		httpClient = http.DefaultClient
	}
	tmpToken := Token{ConsumerKey: token.ConsumerKey, ConsumerKeySecret: token.ConsumerKeySecret}
	return &TwitterEngagementAPI{ctx: ctx, token: tmpToken, httpClient: httpClient, appOnly: true}
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
	header, err := tea.oauthHeader(url)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", header)
	request.Header.Add("Content-Type", "application/json")
	response, err := tea.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
