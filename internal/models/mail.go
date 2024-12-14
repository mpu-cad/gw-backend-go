package models

type Gmail struct {
	Subject     string
	Content     string
	TO          []string
	CC          []string
	BCC         []string
	AttachFiles []string
}
