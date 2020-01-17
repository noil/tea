package tea

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrExceedsMaxTweets    = errors.New("number of tweets exceeds 250")
	ErrExceedsMaxGroupings = errors.New("number of groupings exceeds 3")
	ErrHTTPBadCode         = errors.New("HTTP bad code")
)

type EngagementType string

const (
	TotalImpressionType EngagementType = "impressions" //`90days:true`
	TotalEngagementType EngagementType = "engagements" //`90days:true`
	TotalFavoriteType   EngagementType = "favorites"   //`90days:false`
	TotalRetweetType    EngagementType = "retweets"    //`90days:false`
	TotalReplyType      EngagementType = "replies"     //`90days:false`
	TotalVideoViewType  EngagementType = "video_views" //`90days:false`
)

var defaultEngagementTypes = []EngagementType{
	TotalImpressionType,
	TotalEngagementType,
	TotalFavoriteType,
	TotalRetweetType,
	TotalReplyType,
	TotalVideoViewType,
}

type Grouping map[string]*Group

type Group struct {
	By []By `json:"group_by"`
}

type By string

const (
	TotalTweetIdGroup        By = "tweet.id"
	TotalEngagementTypeGroup By = "engagement.type"
	TotalDefaultGrouping        = "default_grouping"
	UnsupportedTweetIds         = "unsupported_for_impressions_engagements_tweet_ids"
	UnavailableTweetIds         = "unavailable_tweet_ids"
)

var defaultGroupings = Grouping{
	TotalDefaultGrouping: &Group{
		By: []By{
			TotalTweetIdGroup,
			TotalEngagementTypeGroup,
		},
	},
}

type Total struct {
	TweetIds        []string         `json:"tweet_ids"`
	EngagementTypes []EngagementType `json:"engagement_types"`
	Groupings       Grouping         `json:"groupings"`
	tea             *TwitterEngagementAPI
	err             error
	valid           bool
}

func (tea *TwitterEngagementAPI) Total(tweetIds []string) *Total {
	total := &Total{
		tea:   tea,
		valid: true,
	}
	if 250 < len(tweetIds) {
		total.valid = false
		total.err = ErrExceedsMaxTweets
	} else {
		total.TweetIds = tweetIds
	}
	return total
}

func (total *Total) Valid() bool {
	return total.valid
}

func (total *Total) Error() error {
	return total.err
}

func (total *Total) do() (*http.Response, error) {
	if 0 == len(total.EngagementTypes) {
		total.EngagementTypes = defaultEngagementTypes
	}
	if 0 == len(total.Groupings) {
		total.Groupings = defaultGroupings
	}
	params, err := json.Marshal(total)
	if nil != err {
		total.err = err
		total.valid = false
		return nil, err
	}
	request, err := http.NewRequestWithContext(total.tea.ctx, "POST", totalUrl, bytes.NewBuffer(params))
	request.Header.Add("Accept-Encoding", "gzip")
	request.Header.Add("Authorization", total.tea.OAuthHeader(totalUrl))
	request.Header.Add("Content-Type", "application/json")
	response, err := total.tea.httpClient.Do(request)
	if err != nil {
		total.err = err
		total.valid = false
		return nil, err
	}

	return response, nil
}

func (total *Total) Result() (*TotalAPISuccess, error) {
	result := newTotalAPISuccess()
	total.EngagementTypes = defaultEngagementTypes
	total.Groupings = defaultGroupings
	response, err := total.do()
	if nil != err {
		return result, err
	}
	defer func() {
		if nil != response {
			response.Body.Close()
		}
	}()

	result.meta(response)
	reader, err := gzip.NewReader(response.Body)
	if nil != err {
		return nil, err
	}
	if http.StatusOK != response.StatusCode {
		errors := &APIError{}
		err = json.NewDecoder(reader).Decode(errors)
		if nil != err {
			return result, err
		}
		return nil, errors
	}
	tmpSuccess := APISuccessRaw{}
	err = json.NewDecoder(reader).Decode(&tmpSuccess)
	if nil != err {
		return result, err
	}
	result.populate(tmpSuccess)
	return result, nil
}

func (total *Total) ResultRaw(grouping Grouping, types ...EngagementType) (*TotalAPISuccessRaw, error) {
	total.grouping(grouping)
	total.engagementType(types)
	result := newTotalAPISuccessRaw()
	response, err := total.do()
	if nil != err {
		return nil, err
	}
	defer func() {
		if nil != response {
			response.Body.Close()
		}
	}()

	result.meta(response)
	reader, err := gzip.NewReader(response.Body)
	if nil != err {
		return nil, err
	}
	if http.StatusOK != response.StatusCode {
		errors := &APIError{}
		err = json.NewDecoder(reader).Decode(errors)
		if nil != err {
			return nil, err
		}
		return nil, errors
	}
	err = json.NewDecoder(reader).Decode(&result.Data)
	if nil != err {
		return nil, err
	}

	return result, nil
}

func (total *Total) engagementType(types []EngagementType) {
	if 0 == len(types) {
		types = defaultEngagementTypes
	}
	total.EngagementTypes = types
}

func (total *Total) grouping(grouping Grouping) {
	if 0 == len(grouping) {
		grouping = defaultGroupings
	}
	if 3 < len(grouping) {
		total.valid = false
		total.err = ErrExceedsMaxGroupings
	}
	total.Groupings = grouping
}
