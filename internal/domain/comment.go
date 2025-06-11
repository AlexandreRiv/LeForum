package domain

type Comment struct {
	Id         int
	Content    string
	Username   string
	Likes      int
	Dislikes   int
	CreatedAt  string
}