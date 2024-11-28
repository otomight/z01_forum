package config

// TABLES DEFINITIONS
type ClientsTableKeys struct {
	Table			string
	UserId			string
	LastName		string
	FirstName		string
	UserName		string
	Email			string
	OauthProvider	string
	OauthId			string
	Password		string
	Avatar			string
	BirthDate		string
	UserRole		string
	CreationDate	string
	UpdateDate		string
	DeletionDate	string
}

type SessionsTableKeys struct {
	Table			string
	SessionId		string
	UserId			string
	UserRole		string
	UserName		string
	Expiration		string
	CreationDate	string
	UpdateDate		string
	DeletionDate	string
	IsDeleted		string
}

type CategoriesTableKeys struct {
	Table	string
	Id		string
	Name	string
}

type PostsTableKeys struct {
	Table			string
	PostId			string
	AuthorId		string
	Title			string
	Category		string
	Tags			string
	Content			string
	CreationDate	string
	UpdateDate		string
	DeletionDate	string
	IsDeleted		string
	Likes			string
	Dislikes		string
}

type PostsCategoriesTableKeys struct {
	Table		string
	Id			string
	CategoryId	string
	PostId		string
}

type CommentsTableKeys struct {
	Table			string
	CommentId		string
	PostId			string
	UserId			string
	Content			string
	CreationDate	string
}

type LikesDislikesTableKeys struct {
	Table		string
	Id			string
	PostID		string
	UserId		string
	Liked		string
	UpdateDate	string
}

// TABLES IN STRUCT DEFINITION
type StructTablesKeys struct {
	Clients			ClientsTableKeys
	Sessions		SessionsTableKeys
	Categories		CategoriesTableKeys
	Posts			PostsTableKeys
	PostsCategories	PostsCategoriesTableKeys
	Comments		CommentsTableKeys
	LikesDislikes	LikesDislikesTableKeys
}


// TABLES
var clients = ClientsTableKeys{
	Table:			"clients",
	UserId:			"user_id",
	LastName:		"last_name",
	FirstName:		"first_name",
	UserName:		"user_name",
	Email:			"email",
	OauthProvider:	"oauth_provider",
	OauthId:		"oauth_id",
	Password:		"password",
	Avatar:			"avatar",
	BirthDate:		"birth_date",
	UserRole:		"user_role",
	CreationDate:	"creation_date",
	UpdateDate:		"update_date",
	DeletionDate:	"deletion_date",
}

var sessions = SessionsTableKeys{
	Table:			"sessions",
	SessionId:		"session_id",
	UserId:			"user_id",
	UserRole:		"user_role",
	UserName:		"user_name",
	Expiration:		"expiration",
	CreationDate:	"creation_date",
	UpdateDate:		"update_date",
	DeletionDate:	"deletion_date",
	IsDeleted:		"is_deleted",
}

var categories = CategoriesTableKeys{
	Table:	"categories",
	Id:		"id",
	Name:	"name",
}

var posts = PostsTableKeys{
	Table:			"posts",
	PostId:			"post_id",
	AuthorId:		"author_id",
	Title:			"title",
	Category:		"category",
	Tags:			"tags",
	Content:		"content",
	CreationDate:	"creation_date",
	UpdateDate:		"update_date",
	DeletionDate:	"deletion_date",
	IsDeleted:		"is_deleted",
	Likes:			"likes",
	Dislikes:		"dislikes",
}

var postsCategories = PostsCategoriesTableKeys{
	Table:		"posts_categories",
	Id:			"id",
	CategoryId:	"category_id",
	PostId:		"post_id",
}

var comments = CommentsTableKeys{
	Table:			"comments",
	CommentId:		"comment_id",
	PostId:			"post_id",
	UserId:			"user_id",
	Content:		"content",
	CreationDate:	"creation_date",
}

var likesDislikes = LikesDislikesTableKeys{
	Table:		"likes_dislikes",
	Id:			"id",
	PostID:		"post_id",
	UserId:		"user_id",
	Liked:		"liked",
	UpdateDate:	"update_date",
}

// TABLES IN STRUCT
var TableKeys = StructTablesKeys{
	Clients:			clients,
	Sessions:			sessions,
	Categories:			categories,
	Posts:				posts,
	PostsCategories:	postsCategories,
	Comments:			comments,
	LikesDislikes:		likesDislikes,
}


var CategoriesNames = []string{
	"category1",
	"category2",
	"category3",
	"category4",
}

const (
	DbFilePath = "forum.db"
	SqlTablesFilePath = "internal/database/create_table.sql"
)
