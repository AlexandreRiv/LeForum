package domain

type Post struct {
	Id         int
	Title      string
	Content    string
<<<<<<< Updated upstream
	ImageURL   string
	Image      []byte
=======
	Image	   []byte
>>>>>>> Stashed changes
	Username   string
	Categories []string
	Likes      int
	Dislikes   int
	Comments   int
	CreatedAt  string
}
