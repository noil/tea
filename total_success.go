package tea

import (
	"net/http"
	"strconv"
)

type TotalAPISuccess struct {
	Meta   HTTPMeta
	Tweets []Tweet
	Types  Types
}

type Tweet struct {
	ID    string
	Valid bool
	Value int64
	Types Types
}

type Types struct {
	Impression Type
	Engagement Type
	Favorite   Type
	Retweet    Type
	Reply      Type
	VideoView  Type
}

type Type struct {
	Valid bool
	Int64 int64
}

func newTotalAPISuccess() *TotalAPISuccess {
	return &TotalAPISuccess{
		Meta: HTTPMeta{Headers: make(Headers)},
	}
}

func (result *TotalAPISuccess) meta(response *http.Response) {
	if nil != response {
		result.Meta.Code = response.StatusCode
		result.Meta.Headers = headers(response)
	}
}

func (result *TotalAPISuccess) populate(data APISuccessRaw) {
	var (
		types map[string]interface{}
		tweet map[string]interface{}
		ids   []interface{}
		ok    bool
	)
	for label, tweets := range data {
		switch label {
		case TotalDefaultGrouping:
			if tweet, ok = tweets.(map[string]interface{}); !ok {
				continue
			}
			for tweetId, insights := range tweet {
				tweet := Tweet{
					ID:    tweetId,
					Valid: true,
				}
				if types, ok = insights.(map[string]interface{}); !ok {
					continue
				}
				for metric, ivalue := range types {
					value, ok := ivalue.(string)
					if !ok {
						continue
					}
					i, err := strconv.ParseInt(value, 10, 64)
					if nil != err {
						continue
					}
					tweet.Value += i
					switch EngagementType(metric) {
					case TotalImpressionType:
						tweet.Types.Impression.Int64 += i
						tweet.Types.Impression.Valid = true
						result.Types.Impression.Int64 += i
						result.Types.Impression.Valid = true
					case TotalEngagementType:
						tweet.Types.Engagement.Int64 += i
						tweet.Types.Engagement.Valid = true
						result.Types.Engagement.Int64 += i
						result.Types.Engagement.Valid = true
					case TotalFavoriteType:
						tweet.Types.Favorite.Int64 += i
						tweet.Types.Favorite.Valid = true
						result.Types.Favorite.Int64 += i
						result.Types.Favorite.Valid = true
					case TotalRetweetType:
						tweet.Types.Retweet.Int64 += i
						tweet.Types.Retweet.Valid = true
						result.Types.Retweet.Int64 += i
						result.Types.Retweet.Valid = true
					case TotalReplyType:
						tweet.Types.Reply.Int64 += i
						tweet.Types.Reply.Valid = true
						result.Types.Reply.Int64 += i
						result.Types.Reply.Valid = true
					case TotalVideoViewType:
						tweet.Types.VideoView.Int64 += i
						tweet.Types.VideoView.Valid = true
						result.Types.VideoView.Int64 += i
						result.Types.VideoView.Valid = true
					}

				}
				result.Tweets = append(result.Tweets, tweet)
			}
		case UnsupportedTweetIds:
		case UnavailableTweetIds:
			if ids, ok = tweets.([]interface{}); !ok {
				continue
			}
			for _, iid := range ids {
				if id, ok := iid.(string); ok {
					tweet := Tweet{
						ID:    id,
						Valid: false,
					}
					result.Tweets = append(result.Tweets, tweet)
				}
			}
		}
	}
}
