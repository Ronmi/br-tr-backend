package main

import (
	"encoding/json"
	"net/http"

	"github.com/Patrolavia/jsonapi"
)

type api struct {
	data map[string]Project
}

func (a *api) list(dec *json.Decoder, httpData *jsonapi.HTTP) (interface{}, error) {
	ret := make([]Project, 0, len(a.data))
	for _, p := range a.data {
		ret = append(ret, p)
	}
	return ret, nil
}

func (a *api) setOwner(dec *json.Decoder, httpData *jsonapi.HTTP) (interface{}, error) {
	type p struct {
		Repo   string `json:"repo"`
		Branch string `json:"branch"`
		Owner  string `json:"owner"`
	}

	var param p
	if err := dec.Decode(&param); err != nil {
		return nil, jsonapi.Error{http.StatusBadRequest, "Parameter error"}
	}

	found := false
	proj, ok := a.data[param.Repo]
	if !ok {
		return nil, jsonapi.Error{http.StatusNotFound, "Repository not found"}
	}
	for i, b := range proj.Branches {
		if b.Name != param.Branch {
			continue
		}

		found = true
		a.data[param.Repo].Branches[i].Owner = param.Owner
		break
	}

	if !found {
		return nil, jsonapi.Error{http.StatusNotFound, "Branch not found"}
	}

	return nil, nil
}

func (a *api) setDesc(dec *json.Decoder, httpData *jsonapi.HTTP) (interface{}, error) {
	type p struct {
		Repo   string `json:"repo"`
		Branch string `json:"branch"`
		Desc   string `json:"desc"`
	}

	var param p
	if err := dec.Decode(&param); err != nil {
		return nil, jsonapi.Error{http.StatusBadRequest, "Parameter error"}
	}

	found := false
	proj, ok := a.data[param.Repo]
	if !ok {
		return nil, jsonapi.Error{http.StatusNotFound, "Repository not found"}
	}
	for i, b := range proj.Branches {
		if b.Name != param.Branch {
			continue
		}

		found = true
		a.data[param.Repo].Branches[i].Desc = param.Desc
		break
	}

	if !found {
		return nil, jsonapi.Error{http.StatusNotFound, "Branch not found"}
	}

	return nil, nil
}
