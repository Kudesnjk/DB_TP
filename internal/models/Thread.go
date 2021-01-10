package models

import (
	"time"
)

type Thread struct {
	Author    string    `json:"author"`
	Created   time.Time `json:"created"`
	ForumSlug string    `json:"forum"`
	ID        uint64    `json:"id"`
	Message   string    `json:"message"`
	Slug      string    `json:"slug,omitempty"`
	Title     string    `json:"title"`
	Votes     int       `json:"votes"`
}
