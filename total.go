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

type Grouping map[string]*GroupBy

type GroupBy struct {
	Groups []Group `json:"group_by"`
}

type Group string

const (
	TotalTweetIdGroup        Group = "tweet.id"
	TotalEngagementTypeGroup Group = "engagement.type"
)

type Total struct {
	tea             *TwitterEngagementAPI
	TweetIds        []string         `json:"tweet_ids"`
	EngagementTypes []EngagementType `json:"engagement_types"`
	Groupings       Grouping         `json:"groupings"`
	valid           bool
	err             error
}

func (tea *TwitterEngagementAPI) Total(tweetIds []string) *Total {
	total := &Total{tea: tea, valid: true}
	if 250 < len(tweetIds) {
		total.valid = false
		total.err = ErrExceedsMaxTweets
	} else {
		total.TweetIds = tweetIds
	}
	return total
}

func (total *Total) EngagementType(types []EngagementType) *Total {
	total.EngagementTypes = types
	return total
}

func (total *Total) Grouping(grouping Grouping) *Total {
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

func (total *Total) Result() (map[string]interface{}, error) {
	params, err := json.Marshal(total)
	if nil != err {
		return nil, err
	}
	request, err := http.NewRequest("POST", totalUrl, bytes.NewBuffer(params))
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
		return nil, err
	}

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
	success := Success{}
	err = json.NewDecoder(reader).Decode(&success)
	if nil != err {
		return nil, err
	}

	return success, nil
}
