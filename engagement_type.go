package tea

// EngagementType kind of metric type
type EngagementType string

const (
	// ImpressionType engagement type impressions
	ImpressionType EngagementType = "impressions" //`90days:true`
	// EngType engagement type engagements
	EngType EngagementType = "engagements" //`90days:true`
	// FavoriteType engagement type favorites
	FavoriteType EngagementType = "favorites" //`90days:false`
	// RetweetType engagement type retweets
	RetweetType EngagementType = "retweets" //`90days:false`
	// ReplyType engagement type replies
	ReplyType EngagementType = "replies" //`90days:false`
	// VideoViewType engagement type video_views
	VideoViewType EngagementType = "video_views" //`90days:false`
	// MediaViewsType engagement type media_views
	MediaViewsType EngagementType = "media_views" //`90days:false`
	// MediaEngagementsType engagement type media_engagements
	MediaEngagementsType EngagementType = "media_engagements" //`90days:false`
	// URLClicksType engagement type url_clicks
	URLClicksType EngagementType = "url_clicks" //`90days:false`
	// HashTagClicksType engagement type hashtag_clicks
	HashTagClicksType EngagementType = "hashtag_clicks" //`90days:false`
	// DetailExpandsType engagement type detail_expands
	DetailExpandsType EngagementType = "detail_expands" //`90days:false`
	// PermalinkClicksType engagement type permalink_clicks
	PermalinkClicksType EngagementType = "permalink_clicks" //`90days:false`
	// AppInstallAttemptsType engagement type app_install_attempts
	AppInstallAttemptsType EngagementType = "app_install_attempts" //`90days:false`
	// AppOpensType engagement type app_opens
	AppOpensType EngagementType = "app_opens" //`90days:false`
	// EmailTweetType engagement type email_tweet
	EmailTweetType EngagementType = "email_tweet" //`90days:false`
	// UserFollowsType engagement type user_follows
	UserFollowsType EngagementType = "user_follows" //`90days:false`
	// UserProfileClicksType engagement type user_profile_clicks
	UserProfileClicksType EngagementType = "user_profile_clicks" //`90days:false`
)
