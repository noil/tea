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
	TotalReplieType     EngagementType = "replies"     //`90days:false`
	TotalVideoViewType  EngagementType = "video_views" //`90days:false`
)

var defaultEngagementTypes = []EngagementType{
	TotalImpressionType,
	TotalEngagementType,
	TotalFavoriteType,
	TotalRetweetType,
	TotalReplieType,
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
)

var defaultGrouping = Grouping{
	"default_grouping": &Group{
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

func (total *Total) EngagementType(types []EngagementType) *Total {
	if 0 == len(types) {
		types = defaultEngagementTypes
	}
	total.EngagementTypes = types
	return total
}

func (total *Total) Grouping(grouping Grouping) *Total {
	if 0 == len(grouping) {
		grouping = defaultGrouping
	}
	if 3 < len(grouping) {
		total.valid = false
		total.err = ErrExceedsMaxGroupings
	}
	total.Groupings = grouping
	return total
}

func (total *Total) Valid() bool {
	return total.valid
}

func (total *Total) Error() error {
	return total.err
}

func (total *Total) Result() (*TotalResult, error) {
	result := newTotalResult()
	params, err := json.Marshal(total)
	if nil != err {
		return result, err
	}
	request, err := http.NewRequestWithContext(total.tea.ctx, "POST", totalUrl, bytes.NewBuffer(params))
	request.Header.Add("Accept-Encoding", "gzip")
	request.Header.Add("Authorization", total.tea.OAuthHeader(totalUrl))
	request.Header.Add("Content-Type", "application/json")
	response, err := total.tea.httpClient.Do(request)
	defer func() {
		if nil != response {
			response.Body.Close()
		}
	}()
	if err != nil {
		return result, err
	}
	result.meta(response)
	reader, err := gzip.NewReader(response.Body)
	if nil != err {
		return result, err
	}
	if http.StatusOK != response.StatusCode {
		errors := &APIError{}
		err = json.NewDecoder(reader).Decode(errors)
		if nil != err {
			return result, err
		}
		return result, errors
	}
	err = json.NewDecoder(reader).Decode(&result.Data)
	if nil != err {
		return result, err
	}

	return result, nil
}
