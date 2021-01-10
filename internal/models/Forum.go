package models

type Forum struct {
	Slug    string `json:"slug"`
	User    string `json:"user"`
	Title   string `json:"title"`
	Posts   uint32 `json:"posts"`
	Threads uint32 `json:"threads"`
}
