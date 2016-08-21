package main

// DataStore abstracts how we store project and branch info
type DataStore interface {
	UpdateOwner(repo, branch, owner string) error
	UpdateDesc(repo, branch, desc string) error
	List() ([]Project, error)

	// insert or update projects into store
	AddProjects(ps []Project) error
	//AddBranch(repo string, br Branch) error
}
