package domain

import "html/template"

type Post struct {
	Id         int
	Title      string
	Content    string
	Image	   template.URL
	Username   string
	Categories []string
	Likes      int
	Dislikes   int
	Comments   int
	CreatedAt  string
}
