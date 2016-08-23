package main

import (
	"encoding/json"
	"net/http"

	"github.com/Patrolavia/jsonapi"
)

type webhook struct {
	gitlab GitlabProvider
	store  DataStore
}

type webhookProject struct {
	Path string `json:"path_with_namespace"`
}
type webhookAttr struct {
	State           string         `json:"state"`
	SourceProjectID int            `json:"source_project_id"`
	TargetProjectID int            `json:"target_project_id"`
	Source          webhookProject `json:"source,omitempty"`
	Target          webhookProject `json:"target,omitempty"`
}
type webhookPayload struct {
	Kind      string         `json:"object_kind"`
	ProjectID int            `json:"project_id,omitempty"`
	Project   webhookProject `json:"project,omitempty"`
	Attr      webhookAttr    `json:"object_attributes,omitempty"`
}

func (w *webhook) entry(dec *json.Decoder, httpData *jsonapi.HTTP) (ret interface{}, err error) {
	var param webhookPayload
	if err := dec.Decode(&param); err != nil {
		return nil, jsonapi.Error{http.StatusBadRequest, "Incorrect format"}
	}

	id := 0
	name := ""
	switch param.Kind {
	case "push":
		name = param.Project.Path
		id = param.ProjectID
	case "merge_request":
		if param.Attr.State == "merged" {
			name = param.Attr.Target.Path
			id = param.Attr.TargetProjectID
		}
	}

	if name != "" && id > 0 {
		err = w.fetchProject(name, id)
	}
	if err != nil {
		err = jsonapi.Error{http.StatusInternalServerError, err.Error()}
	}

	return
}

func (w *webhook) fetchProject(n string, id int) (err error) {
	brs, err := w.gitlab.Branches(id)
	if err != nil {
		return
	}

	p := Project{
		Name:     n,
		Branches: make([]Branch, 0, len(brs)),
	}

	for _, br := range brs {
		p.Branches = append(p.Branches, Branch{Name: br})
	}

	return w.store.AddProjects([]Project{p})
}
