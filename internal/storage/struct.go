package storage

type Post struct {
    Id         int
    Title      string
    Content    string
    Username   string
    Categories []string
    Likes      int
    Dislikes   int
}