package tea

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"net/http"
)

var (
	// ErrExceedsMaxTweets shows that you reach max count of tweets. (250 max)
	ErrExceedsMaxTweets = errors.New("number of tweets exceeds 250")
	// ErrExceedsMaxGroupings provide that you reach max counf of grouping. (3 max)
	ErrExceedsMaxGroupings = errors.New("number of groupings exceeds 3")
	// ErrHTTPBadCode provide that twitter return error
	ErrHTTPBadCode = errors.New("HTTP bad code")
)

var defaultTotalEngagementTypes = []EngagementType{
	ImpressionType,
	EngType,
	FavoriteType,
	RetweetType,
	ReplyType,
	VideoViewType,
}

const (
	maxTotalTweets = 250
)

var defaultTotalGroupings = Grouping{
	defaultGrouping: &Group{
		By: []GroupByType{
			TweetIDGroup,
			EngagementTypeGroup,
		},
	},
}

// Total struct
type total struct {
	TweetIds        []string         `json:"tweet_ids"`
	EngagementTypes []EngagementType `json:"engagement_types"`
	Groupings       Grouping         `json:"groupings"`
}

// TotalRaw return pointer for struct ResponseRaw
func (tea *TwitterEngagementAPI) TotalRaw(tweetIds []string, types []EngagementType, groups Groups) (*ResponseRaw, error) {
	if maxTotalTweets < len(tweetIds) {
		return nil, ErrExceedsMaxTweets
	}
	total := &total{}
	total.TweetIds = tweetIds
	err := total.groups(groups)
	if err != nil {
		return nil, err
	}
	total.engagementType(types)
	if err != nil {
		return nil, err
	}
	result := newResponseRaw()
	response, err := tea.do(totalURL, total)
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

// Total return pointer for struct Success
func (tea *TwitterEngagementAPI) Total(tweetIds []string) (*Success, error) {
	if maxTotalTweets < len(tweetIds) {
		return nil, ErrExceedsMaxTweets
	}
	total := &total{}
	total.TweetIds = tweetIds
	total.EngagementTypes = defaultTotalEngagementTypes
	total.Groupings = defaultTotalGroupings
	result := newSuccess()
	response, err := tea.do(totalURL, total)
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
	tmpSuccess := SuccessRaw{}
	err = json.NewDecoder(reader).Decode(&tmpSuccess)
	if nil != err {
		return result, err
	}
	result.populate(tmpSuccess)
	return result, nil
}

func (total *total) engagementType(types []EngagementType) error {
	if 0 == len(types) {
		types = defaultTotalEngagementTypes
	}
	total.EngagementTypes = types

	//TODO: implement checking for using valid type for each endpoints
	return nil
}

func (total *total) groups(groups Groups) error {
	if 0 == len(groups) {
		total.Groupings = defaultTotalGroupings
	} else {
		if 3 < len(groups) {
			return ErrExceedsMaxGroupings
		}
		tmpGrouping := make(Grouping)
		for label, values := range groups {
			tmpGrouping[label] = &Group{By: values}
		}

		total.Groupings = tmpGrouping
	}
	return nil
}
