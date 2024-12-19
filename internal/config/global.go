package config

type UserRoleStruct struct {
	User		string
	Moderator	string
	Admin		string
}

const (
	DbFilePath				= "forum.db"
	SqlTablesFilePath		= "forum.sql"
	ServerCertifFilePath	= "server.crt"
	ServerKeyFilePath		= "server.key"

	DataDirPath			= "data/"
	ImagesDirPath		= DataDirPath + "images/"
	PostsImagesDirPath	= ImagesDirPath + "posts/"
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
