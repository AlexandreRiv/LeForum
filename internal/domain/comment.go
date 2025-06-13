package domain

import "html/template"

type Comment struct {
	Id         int
	Content    string
	Image	   template.URL
	Username   string
	Likes      int
	Dislikes   int
	CreatedAt  string
}