package models

type LoginErrorMsg struct {
	UserNotFound		string
	IncorrectPassword	string
}

type RegisterErrorMsg struct {
	UsernameAlreadyTaken	string
	EmailAlreadyTaken		string
}
