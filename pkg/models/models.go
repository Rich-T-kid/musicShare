package models


// Only being used to pass username through context.context throughout the life cycle of the request
type UsernameKey struct{}

type SpotifyAuthError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

/*
Regular Error Object
Apart from the response code, unsuccessful responses return a JSON object containing the following information
*/
type SpotifyError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	Exp         int    `json:"exp"`
	Refresh     string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	Expir       int    `json:"exp"`
}
