package main

import (
	"strconv"

	gogitlab "github.com/plouc/go-gitlab-client"
)

type GogitlabProvider struct {
	Client *gogitlab.Gitlab
}

func (p *GogitlabProvider) Branches(proj int) (ret []string, err error) {
	brs, err := p.Client.ProjectBranches(strconv.Itoa(proj))
	if err != nil {
		return
	}

	ret = make([]string, 0, len(brs))
	for _, br := range brs {
		ret = append(ret, br.Name)
	}

	return
}
