package tea

type EngagementType string

const (
	TypeImpressions EngagementType = "impressions"
	TypeEngagements EngagementType = "engagements"
	TypeFavorites   EngagementType = "favorites"
	TypeRetweets    EngagementType = "retweets"
	TypeReplies     EngagementType = "replies"
	TypeVideoViews  EngagementType = "video_views"
)

type GroupBy string

const (
	GroupByTweetId GroupBy = "tweet.id"
	GroupByEngType GroupBy = "engagement.type"
	// EngDay  GroupBy = "engagement.day"
	// EngHour GroupBy = "engagement.hour"
)

type TwitterEngagementAPI struct {
	Token Token
}

var totalUrl = "https://data-api.twitter.com/insights/engagement/totals"
var eng28hrUrl = "https://data-api.twitter.com/insights/engagement/28hr"
var historicalUrl = "https://data-api.twitter.com/insights/engagement/historical"

type Token struct {
	Access            string
	AccessSecret      string
	ConsumerKey       string
	ConsumerKeySecret string
}

func New(token Token) *TwitterEngagementAPI {
	return &TwitterEngagementAPI{Token: token}
}
