package domain

type Post struct {
	Id         int
	Title      string
	Content    string
	Image	   []byte
	Username   string
	Categories []string
	Likes      int
	Dislikes   int
	Comments   int
	CreatedAt  string
}
