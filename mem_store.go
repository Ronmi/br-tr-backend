package main

import "fmt"

// MemStore stores project and branch info in memory
type MemStore struct {
	projects []Project
}

func (s *MemStore) List() (ret []Project, err error) {
	ret = s.projects
	if ret == nil {
		ret = []Project{}
	}

	return
}

func (s *MemStore) UpdateOwner(repo, branch, owner string) (err error) {
	found := false
	for i, p := range s.projects {
		if p.Name != repo {
			continue
		}

		for j, br := range p.Branches {
			if br.Name != branch {
				continue
			}

			found = true
			s.projects[i].Branches[j].Owner = owner
			break
		}
	}

	if !found {
		err = fmt.Errorf("project %s or branch %s is not found", repo, branch)
	}
	return
}

func (s *MemStore) UpdateDesc(repo, branch, desc string) (err error) {
	found := false
	for i, p := range s.projects {
		if p.Name != repo {
			continue
		}

		for j, br := range p.Branches {
			if br.Name != branch {
				continue
			}

			found = true
			s.projects[i].Branches[j].Desc = desc
			break
		}
	}

	if !found {
		err = fmt.Errorf("project %s or branch %s is not found", repo, branch)
	}
	return
}
