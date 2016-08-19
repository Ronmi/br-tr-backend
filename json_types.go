package main

// Branch represents a branch in JSON format, which is used in API interface
type Branch struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
	Desc  string `json:"desc"`
}

// Project represents a project with all branches in it
type Project struct {
	Name     string   `json:"name"`
	Branches []Branch `json:"branches"`
}
