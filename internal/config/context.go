package config

type contextKey string

// Constants
const (
	SessionKey contextKey = "session"
)

// Oauth 2.0 credentials/URLs
const (
	//Google
	GoogleClientID    = "336521994095-0f7tu06nm9juo9v0rfi08g6c6cdciu67.apps.googleusercontent.com"
	GoogleRedirectURI = "https://localhost/auth/google/callback"
	GoogleAuthURL     = "https://accounts.google.com/o/oauth2/v2/auth"
	GoogleTokenURL    = "https://oauth2.googleapis.com/token"
	GoogleUserInfoURL = "https://www.googleapis.com/oauth2/v3/userinfo"

	//Facebook
	FacebookClientID    = "1115219433947595"
	FacebookRedirectURI = "https://localhost/auth/facebook/callback"
	FacebookAuthURL     = "https://www.facebook.com/v21.0/dialog/oauth"
	FacebookTokenURL    = "https://graph.facebook.com/v21.0/oauth/access_token"
	FacebookUserInfoURL = "https://graph.facebook.com/me"

	//Discord
	DiscordClientID    = "1317915050664792075"
	DiscordRedirectURI = "https://localhost/auth/discord/callback"
	DiscordAuthURL     = "https://discord.com/oauth2/authorize"
	DiscordTokenURL    = "https://discord.com/api/oauth2/token"
	DiscordUserInfoURL = "https://discord.com/api/v10/users/@me"
)

type ProviderConfig struct {
	Name         string
	TokenURL     string
	UserInfoURL  string
	ClientID     string
	ClientSecret string
	RedirectURI  string
}
