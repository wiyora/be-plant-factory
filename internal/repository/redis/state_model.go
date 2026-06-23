package redis

type StatePayload struct {
	State       string `json:"state"`
	RedirectUri string `json:"redirect_uri"`
	BaseURL     string `json:"base_url"`
}
