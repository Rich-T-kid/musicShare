package routes

type signIn struct {
	Username string `json:"username"`
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	Exp         int    `json:"exp"`
	Refresh     string `json:"refresh_token"`
}

type refreshResponse struct {
	AccessToken string `json:"access_token"`
	Expir       int    `json:"exp"`
}
