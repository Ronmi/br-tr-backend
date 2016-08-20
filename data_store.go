package main

// DataStore abstracts how we store project and branch info
type DataStore interface {
	UpdateOwner(repo, branch, owner string) error
	UpdateDesc(repo, branch, desc string) error
	List() ([]Project, error)

	// TODO: We will need theses methods to fully operate on the db,
	//       but right now, the signature is only list here as reminder.
	// AddProject(p Project) error
	// AddBranch(repo string, br Branch) error
}
