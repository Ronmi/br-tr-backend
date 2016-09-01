package main

import (
	"github.com/Ronmi/gitlab"
	"github.com/Ronmi/gitlab/webhook"
)

type wh struct {
	gitlab *gitlab.GitLab
	store  DataStore
	push   chan webhook.PushEvent
	mr     chan webhook.MergeRequestEvent
}

func (w *wh) p() {
	for e := range w.push {
		w.fetchProject(e.Project.PathWithNamespace, e.ProjectID)
	}
}

func (w *wh) m() {
	for e := range w.mr {
		w.fetchProject(e.ObjectAttributes.Target.PathWithNamespace, e.ObjectAttributes.TargetProjectID)
	}
}

func (w *wh) start() {
	go w.p()
	go w.m()
}

func (w *wh) fetchProject(n string, id int) {
	brs, pages, err := w.gitlab.ListBranches(id, nil)
	if err != nil {
		return
	}
	for pages.TotalPages > pages.Page {
		var _b []gitlab.Branch
		_b, pages, err = w.gitlab.ListBranches(id, &gitlab.ListOption{Page: pages.NextPage})
		if err != nil {
			return
		}
		brs = append(brs, _b...)
	}

	p := Project{
		Name:     n,
		Branches: make([]Branch, 0, len(brs)),
	}

	for _, br := range brs {
		p.Branches = append(p.Branches, Branch{Name: br.Name})
	}

	w.store.AddProjects([]Project{p})
}
