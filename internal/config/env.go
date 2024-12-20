package config

import "os"

type EnvVarStruct struct {
	GoogleClientSecret	string
}

var EnvVar = EnvVarStruct{}

func (ev *EnvVarStruct) Set() {
	EnvVar = EnvVarStruct{
		GoogleClientSecret:	os.Getenv("GOOGLE_CLIENT_SECRET"),
	}
}
