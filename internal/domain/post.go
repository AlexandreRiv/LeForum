package domain

type Post struct {
	Id         int
	Title      string
	Content    string
	ImageURL   string
	Username   string
	Categories []string
	Likes      int
	Dislikes   int
	Comments   int
	CreatedAt  string
}
