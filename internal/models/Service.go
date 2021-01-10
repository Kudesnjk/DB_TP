package models

type Service struct {
	UsersNum   uint64 `json:"user"`
	ForumsNum  uint64 `json:"forum"`
	ThreadsNum uint64 `json:"thread"`
	PostsNum   uint64 `json:"post"`
}
