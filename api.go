package main

import (
	"encoding/json"

	"github.com/Patrolavia/jsonapi"
	"github.com/Ronmi/br-tr-backend/session"
	"github.com/Ronmi/gitlab"
)

type api struct {
	store  DataStore
	client *gitlab.GitLab
}

func (a *api) list(dec *json.Decoder, httpData *jsonapi.HTTP) (interface{}, error) {
	return a.store.List()
}

func (a *api) setOwner(dec *json.Decoder, httpData *jsonapi.HTTP) (interface{}, error) {
	type p struct {
		Repo   string `json:"repo"`
		Branch string `json:"branch"`
		Owner  string `json:"owner"`
	}

	var param p
	if err := dec.Decode(&param); err != nil {
		return nil, jsonapi.E401.SetData("Parameter error")
	}

	// non of these fields can be empty
	if param.Repo == "" || param.Branch == "" || param.Owner == "" {
		return nil, jsonapi.E401.SetData("Parameter cannot be empty")
	}

	return nil, a.store.UpdateOwner(param.Repo, param.Branch, param.Owner)
}

func (a *api) setDesc(dec *json.Decoder, httpData *jsonapi.HTTP) (interface{}, error) {
	type p struct {
		Repo   string `json:"repo"`
		Branch string `json:"branch"`
		Desc   string `json:"desc"`
	}

	var param p
	if err := dec.Decode(&param); err != nil {
		return nil, jsonapi.E401.SetData("Parameter error")
	}

	// non of these fields can be empty
	if param.Repo == "" || param.Branch == "" || param.Desc == "" {
		return nil, jsonapi.E401.SetData("Parameter cannot be empty")
	}

	return nil, a.store.UpdateDesc(param.Repo, param.Branch, param.Desc)
}

func (a *api) addBranch(dec *json.Decoder, httpData *jsonapi.HTTP) (interface{}, error) {
	type p struct {
		Repo   string `json:"repo"`
		Branch string `json:"branch"`
		Ref    string `json:"ref,omitempty"`
		Desc   string `json:"desc"`
	}

	var param p
	if err := dec.Decode(&param); err != nil {
		return nil, jsonapi.E401.SetData("Parameter error")
	}

	// non of these fields can be empty
	if param.Repo == "" || param.Branch == "" || param.Desc == "" {
		return nil, jsonapi.E401.SetData("Parameter cannot be empty")
	}

	sess := httpData.Request.Context().Value(session.ContextKey).(session.ContextData)
	if sess.Data.Token == nil || sess.Data.User.ID == 0 {
		return nil, jsonapi.E403.SetData("You must login before creating branch")
	}

	if _, ok := a.store.FindProj(param.Repo); !ok {
		return nil, jsonapi.E404.SetData("The project you specified not found, has it been tracked?")
	}

	// create branch on gitlab
	if _, err := a.client.CreateBranch(param.Repo, param.Branch, param.Ref); err != nil {
		return nil, err
	}

	br := Branch{
		Name:  param.Branch,
		Owner: sess.Data.User.Username,
		Desc:  param.Desc,
	}
	if err := a.store.AddBranch(param.Repo, br); err != nil {
		return nil, err
	}

	return nil, nil
}
