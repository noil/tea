package tea

// Groups map of grouping
type Groups map[string][]GroupByType

// GroupByType kind of grouping type
type GroupByType string

const (
	// TweetIDGroup group by tweet_id
	TweetIDGroup GroupByType = "tweet.id"
	// EngagementTypeGroup group by engagement type
	EngagementTypeGroup GroupByType = "engagement.type"
	// EngagementDayGroup group by engagement day
	EngagementDayGroup GroupByType = "engagement.day"
	// EngagementHourGroup group by engagement hour
	EngagementHourGroup GroupByType = "engagement.hour"
)

// Grouping struct
type Grouping map[string]*Group

// Group struct
type Group struct {
	By []GroupByType `json:"group_by"`
}
