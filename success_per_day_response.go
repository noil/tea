package tea

import (
	"net/http"
	"strconv"
	"time"
)

// SuccessPerDay struct SuccessPerDay
type SuccessPerDay struct {
	Meta   Meta
	Tweets []TweetPerDay
	// TODO: implement TypesPerDay
	// Types  Types
}

// TweetPerDay struct TweetPerDay
type TweetPerDay struct {
	ID    string
	Valid bool
	Value int64
	Types Types
	Days  map[string]Types
}

func newSuccessPerDay() *SuccessPerDay {
	return &SuccessPerDay{
		Meta: Meta{HTTP: HTTP{Headers: make(map[string]string)}},
	}
}

func (result *SuccessPerDay) meta(response *http.Response) {
	if nil != response {
		result.Meta.HTTP.Code = response.StatusCode
		result.Meta.HTTP.Headers = headers(response)
	}
}

func (result *SuccessPerDay) populate(data SuccessRaw) {
	var (
		types map[string]interface{}
		tweet map[string]interface{}
		days  map[string]interface{}
		ids   []interface{}
		ok    bool
	)
	for label, tweets := range data {
		switch label {
		case defaultGrouping:
			if tweet, ok = tweets.(map[string]interface{}); !ok {
				continue
			}
			for tweetID, itypes := range tweet {
				tweet := TweetPerDay{
					ID:    tweetID,
					Valid: true,
					Days:  make(map[string]Types),
				}
				if types, ok = itypes.(map[string]interface{}); !ok {
					continue
				}
				for metric, idays := range types {
					if days, ok = idays.(map[string]interface{}); !ok {
						continue
					}
					for day, ivalue := range days {
						tweet.Days[day] = Types{}
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
							dayType := tweet.Days[day]
							dayType.Impression.Int64 += i
							dayType.Impression.Valid = true
							tweet.Days[day] = dayType
						case EngType:
							tweet.Types.Engagement.Int64 += i
							tweet.Types.Engagement.Valid = true
							dayType := tweet.Days[day]
							dayType.Engagement.Int64 += i
							dayType.Engagement.Valid = true
							tweet.Days[day] = dayType
						case FavoriteType:
							tweet.Types.Favorite.Int64 += i
							tweet.Types.Favorite.Valid = true
							dayType := tweet.Days[day]
							dayType.Favorite.Int64 += i
							dayType.Favorite.Valid = true
							tweet.Days[day] = dayType
						case RetweetType:
							tweet.Types.Retweet.Int64 += i
							tweet.Types.Retweet.Valid = true
							dayType := tweet.Days[day]
							dayType.Retweet.Int64 += i
							dayType.Retweet.Valid = true
							tweet.Days[day] = dayType
						case ReplyType:
							tweet.Types.Reply.Int64 += i
							tweet.Types.Reply.Valid = true
							dayType := tweet.Days[day]
							dayType.Reply.Int64 += i
							dayType.Reply.Valid = true
							tweet.Days[day] = dayType
						case VideoViewType:
							tweet.Types.VideoView.Int64 += i
							tweet.Types.VideoView.Valid = true
							dayType := tweet.Days[day]
							dayType.VideoView.Int64 += i
							dayType.VideoView.Valid = true
							tweet.Days[day] = dayType
						case MediaViewsType:
							tweet.Types.MediaView.Int64 += i
							tweet.Types.MediaView.Valid = true
							dayType := tweet.Days[day]
							dayType.MediaView.Int64 += i
							dayType.MediaView.Valid = true
							tweet.Days[day] = dayType
						case MediaEngagementsType:
							tweet.Types.MediaEngagement.Int64 += i
							tweet.Types.MediaEngagement.Valid = true
							dayType := tweet.Days[day]
							dayType.MediaEngagement.Int64 += i
							dayType.MediaEngagement.Valid = true
							tweet.Days[day] = dayType
						case URLClicksType:
							tweet.Types.URLClick.Int64 += i
							tweet.Types.URLClick.Valid = true
							dayType := tweet.Days[day]
							dayType.URLClick.Int64 += i
							dayType.URLClick.Valid = true
							tweet.Days[day] = dayType
						case HashTagClicksType:
							tweet.Types.HashtagClick.Int64 += i
							tweet.Types.HashtagClick.Valid = true
							dayType := tweet.Days[day]
							dayType.HashtagClick.Int64 += i
							dayType.HashtagClick.Valid = true
							tweet.Days[day] = dayType
						case DetailExpandsType:
							tweet.Types.DetailClick.Int64 += i
							tweet.Types.DetailClick.Valid = true
							dayType := tweet.Days[day]
							dayType.DetailClick.Int64 += i
							dayType.DetailClick.Valid = true
							tweet.Days[day] = dayType
						case PermalinkClicksType:
							tweet.Types.PermalinkClick.Int64 += i
							tweet.Types.PermalinkClick.Valid = true
							dayType := tweet.Days[day]
							dayType.PermalinkClick.Int64 += i
							dayType.PermalinkClick.Valid = true
							tweet.Days[day] = dayType
						case AppInstallAttemptsType:
							tweet.Types.AppInstallAttempt.Int64 += i
							tweet.Types.AppInstallAttempt.Valid = true
							dayType := tweet.Days[day]
							dayType.AppInstallAttempt.Int64 += i
							dayType.AppInstallAttempt.Valid = true
							tweet.Days[day] = dayType
						case AppOpensType:
							tweet.Types.AppOpen.Int64 += i
							tweet.Types.AppOpen.Valid = true
							dayType := tweet.Days[day]
							dayType.AppOpen.Int64 += i
							dayType.AppOpen.Valid = true
							tweet.Days[day] = dayType
						case EmailTweetType:
							tweet.Types.TweetEmail.Int64 += i
							tweet.Types.TweetEmail.Valid = true
							dayType := tweet.Days[day]
							dayType.TweetEmail.Int64 += i
							dayType.TweetEmail.Valid = true
							tweet.Days[day] = dayType
						case UserFollowsType:
							tweet.Types.UserFollow.Int64 += i
							tweet.Types.UserFollow.Valid = true
							dayType := tweet.Days[day]
							dayType.UserFollow.Int64 += i
							dayType.UserFollow.Valid = true
							tweet.Days[day] = dayType
						case UserProfileClicksType:
							tweet.Types.UserProfileClick.Int64 += i
							tweet.Types.UserProfileClick.Valid = true
							dayType := tweet.Days[day]
							dayType.UserProfileClick.Int64 += i
							dayType.UserProfileClick.Valid = true
							tweet.Days[day] = dayType
						}
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
					tweet := TweetPerDay{
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
			dt, err := time.Parse("2006-01-02T15:04:05Z", tmpDt)
			if nil != err {
				result.Meta.End = Duration{Valid: false, Time: dt}
			}
			result.Meta.End = Duration{Valid: true, Time: dt}
		}
	}
}
