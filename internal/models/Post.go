package models

import (
	"time"
)

type Post struct {
	Author     string    `json:"author"`
	Created    time.Time `json:"created"`
	ForumSlug  string    `json:"forum"`
	ID         uint64    `json:"id"`
	IsEdited   bool      `json:"isEdited"`
	Message    string    `json:"message"`
	Parent     uint64    `json:"parent"`
	ThreadID   uint64    `json:"thread"`
	Path       []int64   `json:"-"`
	ThreadSlug string    `json:"-"`
}

type AdditionalPostData struct {
	Created    time.Time
	ThreadSlug string
	ForumSlug  string
	ThreadID   uint64
}
