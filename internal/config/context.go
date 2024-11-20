package config

type contextKey string

// Constants
const (
	SessionKey	contextKey = "session"
)

// Oauth 2.0 credentials/URLs
const (
	//Google
	GoogleClientID     = "336521994095-0f7tu06nm9juo9v0rfi08g6c6cdciu67.apps.googleusercontent.com"
	GoogleClientSecret = "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"
	GoogleRedirectURI  = "http://localhost:8081/auth/callback?provider=google"
	GoogleAuthURL      = "https://accounts.google.com/o/oauth2/v2/auth"
	GoogleTokenURL     = "https://oauth2.googleapis.com/token"
	GoogleUserInfoURL  = "https://www.googleapis.com/oauth2/v3/userinfo"

	//Github
	GithubClientID     = "Iv23lixCMoiAieKSPEuk"
	GithubClientSecret = "1095fbed2fe4f4790a45bf79986b4543940de3bb"
	GithubRedirectURI  = "http://localhost:8081/auth/github/callback"
	GithubAuthURL      = "https://github.com/login/oauth/authorize"
	GithubTokenURL     = "https://github.com/login/oauth/access_token"
	GithubUserInfoURL  = "https://api.github.com/user"

	//Facebook
	FacebookClientID     = "3907245556220479"
	FacebookClientSecret = "9f945eb300532f82baa8d62d5a613d3e"
	FacebookRedirectURI  = "http://localhost:8081"
	FacebookAuthURL      = "https://www.facebook.com/v12.0/dialog/oauth"
	FacebookTokenURL     = "https://graph.facebook.com/v12.0/oauth/access_token"
	FacebookUserInfoURL  = "https://graph.facebook.com/me?fields=id,name,email"
)

type ProviderConfig struct {
	Name         string
	TokenURL     string
	UserInfoURL  string
	ClientID     string
	ClientSecret string
	RedirectURI  string
}
