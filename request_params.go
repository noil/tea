package tea

import (
	"encoding/json"
)

type reqParams struct {
	TweetIds []string         `json:"tweet_ids"`
	Types    []EngagementType `json:"engagement_types"`
	GroupBy  *gName           `json:"groupings"`
}

type gName struct {
	GBy *gBy `json:"grouping_name"`
}

type gBy struct {
	By []GroupBy `json:"group_by"`
}

func requestParams(tweetIds []string, types []EngagementType, groupBy []GroupBy) ([]byte, error) {
	rp := &reqParams{TweetIds: tweetIds, Types: types,
		GroupBy: &gName{
			GBy: &gBy{
				By: groupBy,
			},
		},
	}
	return json.Marshal(rp)
}
