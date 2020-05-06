package tea

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

const (
	maxEng28Tweets = 25
)

var (
	// ErrExceedsMaxEng28Tweets provide that you reach max count of tweets. (25 max)
	ErrExceedsMaxEng28Tweets = errors.New("number of tweets exceeds 25")
)

var defaultEng28EngagementTypes = []EngagementType{
	ImpressionType,
	EngType,
	FavoriteType,
	RetweetType,
	ReplyType,
	VideoViewType,
	MediaViewsType,
	MediaEngagementsType,
	URLClicksType,
	HashTagClicksType,
	DetailExpandsType,
	PermalinkClicksType,
	AppInstallAttemptsType,
	AppOpensType,
	EmailTweetType,
	UserFollowsType,
	UserProfileClicksType,
}

var defaultEng28Groupings = Grouping{
	defaultGrouping: &Group{
		By: []GroupByType{
			TweetIDGroup,
			EngagementTypeGroup,
			// EngagementDayGroup,
		},
	},
}

// Eng28 struct Eng28
type eng28 struct {
	TweetIds        []string         `json:"tweet_ids"`
	EngagementTypes []EngagementType `json:"engagement_types"`
	Groupings       Grouping         `json:"groupings"`
}

// Eng28Raw return pointer for struct ResponseRaw
func (tea *TwitterEngagementAPI) Eng28Raw(tweetIds []string, types []EngagementType, groups Groups, from, to time.Time) (*ResponseRaw, error) {
	if maxEng28Tweets < len(tweetIds) {
		return nil, ErrExceedsMaxEng28Tweets
	}
	eng28 := &eng28{}
	eng28.TweetIds = tweetIds
	err := eng28.groups(groups)
	if err != nil {
		return nil, err
	}
	eng28.engagementType(types)
	if err != nil {
		return nil, err
	}
	result := newResponseRaw()
	response, err := tea.do(eng28hrURL, eng28)
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
		errors := &EngAPIError{}
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

// Eng28 return pointer for struct Success
func (tea *TwitterEngagementAPI) Eng28(tweetIds []string, from, to time.Time) (*Success, error) {
	if maxEng28Tweets < len(tweetIds) {
		return nil, ErrExceedsMaxEng28Tweets
	}
	eng28 := &eng28{}
	eng28.TweetIds = tweetIds
	eng28.EngagementTypes = defaultEng28EngagementTypes
	eng28.Groupings = defaultEng28Groupings
	result := newSuccess()
	response, err := tea.do(eng28hrURL, eng28)
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
		errors := &EngAPIError{}
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

func (eng28 *eng28) engagementType(types []EngagementType) error {
	if 0 == len(types) {
		types = defaultEng28EngagementTypes
	}
	eng28.EngagementTypes = types

	//TODO: implement checking for using valid type for each endpoints
	return nil
}

func (eng28 *eng28) groups(groups Groups) error {
	if 0 == len(groups) {
		eng28.Groupings = defaultEng28Groupings
	} else {
		if 3 < len(groups) {
			return ErrExceedsMaxGroupings
		}
		tmpGrouping := make(Grouping)
		for label, values := range groups {
			tmpGrouping[label] = &Group{By: values}
		}

		eng28.Groupings = tmpGrouping
	}
	return nil
}
