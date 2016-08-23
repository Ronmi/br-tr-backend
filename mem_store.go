package main

import (
	"fmt"
	"sync"
)

// MemStore stores project and branch info in memory
type MemStore struct {
	projects []Project
	lock     *sync.RWMutex
}

func (s *MemStore) List() (ret []Project, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	ret = s.projects
	if ret == nil {
		ret = []Project{}
	}

	return
}

func (s *MemStore) UpdateOwner(repo, branch, owner string) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
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
	s.lock.Lock()
	defer s.lock.Unlock()
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

func (s *MemStore) AddProjects(ps []Project) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	// create a repo_name => array index map
	lookup := make(map[string]int)
	for idx, p := range s.projects {
		lookup[p.Name] = idx
	}

	// insert or update it
	for _, p := range ps {
		idx, ok := lookup[p.Name]
		if !ok {
			// insert
			lookup[p.Name] = len(s.projects)
			s.projects = append(s.projects, p)
			continue
		}

		// update
		s.projects[idx] = p
	}

	return nil
}
