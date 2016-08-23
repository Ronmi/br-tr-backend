package main

// GitlabProvider is a translating layer between Gitlab API and our need
type GitlabProvider interface {
	Branches(proj int) ([]string, error) // return only branch names
}
