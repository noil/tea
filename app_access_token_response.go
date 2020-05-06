package tea

type appOAuthToken struct {
	Type  string `json:"token_type"`
	Value string `json:"access_token"`
}
