package config

import "os"

type EnvVarStruct struct {
	GoogleClientSecret		string
	FacebookClientSecret	string
	DiscordClientSecret		string
}

var EnvVar = EnvVarStruct{}

func (ev *EnvVarStruct) Set() {
	EnvVar = EnvVarStruct{
		GoogleClientSecret:		os.Getenv("GOOGLE_CLIENT_SECRET"),
		FacebookClientSecret:	os.Getenv("FACEBOOK_CLIENT_SECRET"),
		DiscordClientSecret:	os.Getenv("DISCORD_CLIENT_SECRET"),
	}
}
