package config

// TABLES DEFINITIONS
type ClientsTableKeys struct {
	Clients			string
	ID				string
	LastName		string
	FirstName		string
	UserName		string
	Email			string
	OauthProvider	string
	OauthID			string
	Password		string
	Avatar			string
	BirthDate		string
	UserRole		string
	CreationDate	string
	UpdateDate		string
	DeletionDate	string
}

type SessionsTableKeys struct {
	Sessions		string
	ID				string
	UserID			string
	UserRole		string
	UserName		string
	Expiration		string
	CreationDate	string
	UpdateDate		string
	DeletionDate	string
	IsDeleted		string
}

type CategoriesTableKeys struct {
	Categories	string
	ID			string
	Name		string
}

type PostsTableKeys struct {
	Posts			string
	ID				string
	AuthorID		string
	Title			string
	Content			string
	CreationDate	string
	UpdateDate		string
	DeletionDate	string
	IsDeleted		string
	Likes			string
	Dislikes		string
}

type PostsCategoriesTableKeys struct {
	PostsCategories	string
	ID				string
	CategoryID		string
	PostID			string
}

type CommentsTableKeys struct {
	Comments		string
	ID				string
	PostID			string
	UserID			string
	Content			string
	CreationDate	string
}

type ReactionsTableKeys struct {
	Reactions	string
	ID				string
	PostID			string
	UserID			string
	Liked			string
	UpdateDate		string
}

// TABLES IN STRUCT DEFINITION
type StructTablesKeys struct {
	Clients			ClientsTableKeys
	Sessions		SessionsTableKeys
	Categories		CategoriesTableKeys
	Posts			PostsTableKeys
	PostsCategories	PostsCategoriesTableKeys
	Comments		CommentsTableKeys
	Reactions	ReactionsTableKeys
}


// TABLES
var clients = ClientsTableKeys{
	Clients:		"clients",
	ID:				"id",
	LastName:		"last_name",
	FirstName:		"first_name",
	UserName:		"user_name",
	Email:			"email",
	OauthProvider:	"oauth_provider",
	OauthID:		"oauth_id",
	Password:		"password",
	Avatar:			"avatar",
	BirthDate:		"birth_date",
	UserRole:		"user_role",
	CreationDate:	"creation_date",
	UpdateDate:		"update_date",
	DeletionDate:	"deletion_date",
}

var sessions = SessionsTableKeys{
	Sessions:		"sessions",
	ID:				"id",
	UserID:			"user_id",
	UserRole:		"user_role",
	UserName:		"user_name",
	Expiration:		"expiration",
	CreationDate:	"creation_date",
	UpdateDate:		"update_date",
	DeletionDate:	"deletion_date",
	IsDeleted:		"is_deleted",
}

var categories = CategoriesTableKeys{
	Categories:	"categories",
	ID:			"id",
	Name:		"name",
}

var posts = PostsTableKeys{
	Posts:			"posts",
	ID:				"id",
	AuthorID:		"author_id",
	Title:			"title",
	Content:		"content",
	CreationDate:	"creation_date",
	UpdateDate:		"update_date",
	DeletionDate:	"deletion_date",
	IsDeleted:		"is_deleted",
	Likes:			"likes",
	Dislikes:		"dislikes",
}

var postsCategories = PostsCategoriesTableKeys{
	PostsCategories:	"posts_categories",
	ID:					"id",
	CategoryID:			"category_id",
	PostID:				"post_id",
}

var comments = CommentsTableKeys{
	Comments:		"comments",
	ID:				"id",
	PostID:			"post_id",
	UserID:			"user_id",
	Content:		"content",
	CreationDate:	"creation_date",
}

var reactions = ReactionsTableKeys{
	Reactions:		"reactions",
	ID:				"id",
	PostID:			"post_id",
	UserID:			"user_id",
	Liked:			"liked",
	UpdateDate:		"update_date",
}

// TABLES IN STRUCT
var TableKeys = StructTablesKeys{
	Clients:			clients,
	Sessions:			sessions,
	Categories:			categories,
	Posts:				posts,
	PostsCategories:	postsCategories,
	Comments:			comments,
	Reactions:			reactions,
}


var CategoriesNames = []string{
	"category1",
	"category2",
	"category3",
	"category4",
}

const (
	DbFilePath			= "forum.db"
	SqlTablesFilePath	= "internal/database/create_table.sql"
)
