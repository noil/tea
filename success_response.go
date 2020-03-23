package tea

import (
	"net/http"
	"strconv"
	"time"
)

// Success struct Success
type Success struct {
	Meta   Meta
	Tweets []Tweet
	Types  Types
}

// Tweet struct Tweet
type Tweet struct {
	ID    string
	Valid bool
	Value int64
	Types Types
}

// Types struct Types
type Types struct {
	Impression        Type
	Engagement        Type
	Favorite          Type
	Retweet           Type
	Reply             Type
	VideoView         Type
	MediaView         Type
	MediaEngagement   Type
	URLClick          Type
	HashtagClick      Type
	DetailClick       Type
	PermalinkClick    Type
	AppInstallAttempt Type
	AppOpen           Type
	TweetEmail        Type
	UserFollow        Type
	UserProfileClick  Type
}

// Type struct Type
type Type struct {
	Valid bool
	Int64 int64
}

func newSuccess() *Success {
	return &Success{
		Meta: Meta{HTTP: HTTP{Headers: make(map[string]string)}},
	}
}

func (result *Success) meta(response *http.Response) {
	if nil != response {
		result.Meta.HTTP.Code = response.StatusCode
		result.Meta.HTTP.Headers = headers(response)
	}
}

func (result *Success) populate(data SuccessRaw) {
	var (
		types map[string]interface{}
		tweet map[string]interface{}
		ids   []interface{}
		ok    bool
	)
	for label, tweets := range data {
		switch label {
		case defaultGrouping:
			if tweet, ok = tweets.(map[string]interface{}); !ok {
				continue
			}
			for tweetID, insights := range tweet {
				tweet := Tweet{
					ID:    tweetID,
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
					case ImpressionType:
						tweet.Types.Impression.Int64 += i
						tweet.Types.Impression.Valid = true
						result.Types.Impression.Int64 += i
						result.Types.Impression.Valid = true
					case EngType:
						tweet.Types.Engagement.Int64 += i
						tweet.Types.Engagement.Valid = true
						result.Types.Engagement.Int64 += i
						result.Types.Engagement.Valid = true
					case FavoriteType:
						tweet.Types.Favorite.Int64 += i
						tweet.Types.Favorite.Valid = true
						result.Types.Favorite.Int64 += i
						result.Types.Favorite.Valid = true
					case RetweetType:
						tweet.Types.Retweet.Int64 += i
						tweet.Types.Retweet.Valid = true
						result.Types.Retweet.Int64 += i
						result.Types.Retweet.Valid = true
					case ReplyType:
						tweet.Types.Reply.Int64 += i
						tweet.Types.Reply.Valid = true
						result.Types.Reply.Int64 += i
						result.Types.Reply.Valid = true
					case VideoViewType:
						tweet.Types.VideoView.Int64 += i
						tweet.Types.VideoView.Valid = true
						result.Types.VideoView.Int64 += i
						result.Types.VideoView.Valid = true
					case MediaViewsType:
						tweet.Types.MediaView.Int64 += i
						tweet.Types.MediaView.Valid = true
						result.Types.MediaView.Int64 += i
						result.Types.MediaView.Valid = true
					case MediaEngagementsType:
						tweet.Types.MediaEngagement.Int64 += i
						tweet.Types.MediaEngagement.Valid = true
						result.Types.MediaEngagement.Int64 += i
						result.Types.MediaEngagement.Valid = true
					case URLClicksType:
						tweet.Types.URLClick.Int64 += i
						tweet.Types.URLClick.Valid = true
						result.Types.URLClick.Int64 += i
						result.Types.URLClick.Valid = true
					case HashTagClicksType:
						tweet.Types.HashtagClick.Int64 += i
						tweet.Types.HashtagClick.Valid = true
						result.Types.HashtagClick.Int64 += i
						result.Types.HashtagClick.Valid = true
					case DetailExpandsType:
						tweet.Types.DetailClick.Int64 += i
						tweet.Types.DetailClick.Valid = true
						result.Types.DetailClick.Int64 += i
						result.Types.DetailClick.Valid = true
					case PermalinkClicksType:
						tweet.Types.PermalinkClick.Int64 += i
						tweet.Types.PermalinkClick.Valid = true
						result.Types.PermalinkClick.Int64 += i
						result.Types.PermalinkClick.Valid = true
					case AppInstallAttemptsType:
						tweet.Types.AppInstallAttempt.Int64 += i
						tweet.Types.AppInstallAttempt.Valid = true
						result.Types.AppInstallAttempt.Int64 += i
						result.Types.AppInstallAttempt.Valid = true
					case AppOpensType:
						tweet.Types.AppOpen.Int64 += i
						tweet.Types.AppOpen.Valid = true
						result.Types.AppOpen.Int64 += i
						result.Types.AppOpen.Valid = true
					case EmailTweetType:
						tweet.Types.TweetEmail.Int64 += i
						tweet.Types.TweetEmail.Valid = true
						result.Types.TweetEmail.Int64 += i
						result.Types.TweetEmail.Valid = true
					case UserFollowsType:
						tweet.Types.UserFollow.Int64 += i
						tweet.Types.UserFollow.Valid = true
						result.Types.UserFollow.Int64 += i
						result.Types.UserFollow.Valid = true
					case UserProfileClicksType:
						tweet.Types.UserProfileClick.Int64 += i
						tweet.Types.UserProfileClick.Valid = true
						result.Types.UserProfileClick.Int64 += i
						result.Types.UserProfileClick.Valid = true
					}
				}
				result.Tweets = append(result.Tweets, tweet)
			}
		case unsupportedTweetIds:
		case unavailableTweetIds:
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
		case start:
			var tmpDt string
			if tmpDt, ok = tweets.(string); !ok {
				continue
			}
			dt, err := time.Parse("2006-01-02T15:04:05Z", tmpDt)
			if nil != err {
				result.Meta.Start = Duration{Valid: false, Time: dt}
			}
			result.Meta.Start = Duration{Valid: true, Time: dt}
		case end:
			var tmpDt string
			if tmpDt, ok = tweets.(string); !ok {
				continue
			}
			dt, err := time.Parse("2016-01-02T15:04:05Z", tmpDt)
			if nil != err {
				result.Meta.End = Duration{Valid: false, Time: dt}
			}
			result.Meta.End = Duration{Valid: true, Time: dt}
		}
	}
}
