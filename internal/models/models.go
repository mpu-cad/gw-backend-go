package models

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	LastName string `json:"last_name"`
	Login    string `json:"login"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	HashPass string `json:"hash_pass"`
	IsAdmin  bool   `json:"is_admin"`
	IsBanned bool   `json:"is_banned"`
}

type Gmail struct {
	Subject     string
	Content     string
	TO          []string
	CC          []string
	BCC         []string
	AttachFiles []string
}
