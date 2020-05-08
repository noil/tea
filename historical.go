package tea

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	maxHistoricalTweets = 25
)

var (
	// ErrExceedsMaxHistoricalTweets provide that you reach max count of tweets. (25 max)
	ErrExceedsMaxHistoricalTweets = errors.New("number of tweets exceeds 25")
)

var defaultHistoricalEngagementTypes = []EngagementType{
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

var (
	defaultHistoricalGroupings = Grouping{
		defaultGrouping: &Group{
			By: []GroupByType{
				TweetIDGroup,
				EngagementTypeGroup,
			},
		},
	}

	defaultHistoricalPerDayGroupings = Grouping{
		defaultGrouping: &Group{
			By: []GroupByType{
				TweetIDGroup,
				EngagementTypeGroup,
				EngagementDayGroup,
			},
		},
	}
)

// Historical struct Historical
type historical struct {
	TweetIds        []string         `json:"tweet_ids"`
	EngagementTypes []EngagementType `json:"engagement_types"`
	Groupings       Grouping         `json:"groupings"`
	Start           string           `json:"start"`
	End             string           `json:"end"`
}

// HistoricalRaw return pointer for struct ResponseRaw
func (tea *TwitterEngagementAPI) HistoricalRaw(tweetIds []string, types []EngagementType, groups Groups, from, to time.Time) (*ResponseRaw, error) {
	if maxHistoricalTweets < len(tweetIds) {
		return nil, ErrExceedsMaxHistoricalTweets
	}
	historical := &historical{}
	historical.TweetIds = tweetIds
	historical.Start = from.Format("2006-01-02")
	historical.End = to.Format("2006-01-02")
	err := historical.groups(groups)
	if err != nil {
		return nil, err
	}
	historical.engagementType(types)
	if err != nil {
		return nil, err
	}
	result := newResponseRaw()
	response, err := tea.do(historicalURL, historical)
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

// Historical return pointer for struct Success
func (tea *TwitterEngagementAPI) Historical(tweetIds []string, from, to time.Time) (*Success, error) {
	if maxHistoricalTweets < len(tweetIds) {
		return nil, ErrExceedsMaxHistoricalTweets
	}
	historical := &historical{}
	historical.TweetIds = tweetIds
	historical.Start = from.Format("2006-01-02")
	historical.End = to.Format("2006-01-02")
	historical.EngagementTypes = defaultHistoricalEngagementTypes
	historical.Groupings = defaultHistoricalGroupings
	result := newSuccess()
	response, err := tea.do(historicalURL, historical)
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

// HistoricalPerDay return pointer for struct SuccessPerDay
func (tea *TwitterEngagementAPI) HistoricalPerDay(tweetIds []string, from, to time.Time) (*SuccessPerDay, error) {
	if maxHistoricalTweets < len(tweetIds) {
		return nil, ErrExceedsMaxHistoricalTweets
	}
	historical := &historical{}
	historical.TweetIds = tweetIds
	historical.Start = from.Format("2006-01-02")
	historical.End = to.Format("2006-01-02")
	historical.EngagementTypes = defaultHistoricalEngagementTypes
	historical.Groupings = defaultHistoricalPerDayGroupings
	result := newSuccessPerDay()
	response, err := tea.do(historicalURL, historical)
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
	b, _ := ioutil.ReadAll(reader)
	fmt.Println(string(b))
	os.Exit(1)
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

func (historical *historical) engagementType(types []EngagementType) error {
	if 0 == len(types) {
		types = defaultHistoricalEngagementTypes
	}
	historical.EngagementTypes = types

	//TODO: implement checking for using valid type for each endpoints
	return nil
}

func (historical *historical) groups(groups Groups) error {
	if 0 == len(groups) {
		historical.Groupings = defaultHistoricalGroupings
	} else {
		if 3 < len(groups) {
			return ErrExceedsMaxGroupings
		}
		tmpGrouping := make(Grouping)
		for label, values := range groups {
			tmpGrouping[label] = &Group{By: values}
		}

		historical.Groupings = tmpGrouping
	}
	return nil
}
