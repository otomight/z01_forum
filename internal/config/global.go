package config

type UserRoleStruct struct {
	User		string
	Moderator	string
	Admin		string
}

const (
	DbFilePath			= "forum.db"
	SqlTablesFilePath	= "forum.sql"
)

var	UserRole = UserRoleStruct{
	User:		"user",
	Moderator:	"moderator",
	Admin:		"administrator",
}

var CategoriesNames = []string{
	"category1",
	"category2",
	"category3",
	"category4",
}
