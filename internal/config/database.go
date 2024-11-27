package config

// Tables definitions
type ClientsTable struct {
	Name	string
}

type SessionsTable struct {
	Name	string
}

type PostsTable struct {
	Name	string
}

type CommentsTable struct {
	Name	string
}

type LikesDislikesTable struct {
	Name	string
}

type Tables struct {
	Clients			ClientsTable
	Sessions		SessionsTable
	Posts			PostsTable
	Comments		CommentsTable
	LikesDislikes	LikesDislikesTable
}


// tables assignations
var clients = ClientsTable{
	Name:	"clients",
}

var sessions = SessionsTable{
	Name:	"sessions",
}

var posts = PostsTable{
	Name:	"posts",
}

var comments = CommentsTable{
	Name:	"comments",
}

var likesDislikes = LikesDislikesTable{
	Name:	"likes_dislikes",
}

var Table = Tables{
	Clients:		clients,
	Sessions:		sessions,
	Posts:			posts,
	Comments:		comments,
	LikesDislikes:	likesDislikes,
}

const (
	DbFilePath = "forum.db"
	SqlTablesFilePath = "internal/database/create_table.sql"
)
