package main

import (
	"encoding/json"
	"net/http"

	"github.com/Patrolavia/jsonapi"
)

type api struct {
	store DataStore
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
		return nil, jsonapi.Error{http.StatusBadRequest, "Parameter error"}
	}

	// non of these fields can be empty
	if param.Repo == "" || param.Branch == "" || param.Owner == "" {
		return nil, jsonapi.Error{http.StatusBadRequest, "Parameter cannot be empty"}
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
		return nil, jsonapi.Error{http.StatusBadRequest, "Parameter error"}
	}

	// non of these fields can be empty
	if param.Repo == "" || param.Branch == "" || param.Desc == "" {
		return nil, jsonapi.Error{http.StatusBadRequest, "Parameter cannot be empty"}
	}

	return nil, a.store.UpdateDesc(param.Repo, param.Branch, param.Desc)
}
