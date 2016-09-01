package main

import (
	"encoding/json"
	"testing"

	"github.com/Patrolavia/jsonapi"
)

func makeAPI() *api {
	return &api{makeMemStore(), nil}
}

func TestAPIList(t *testing.T) {
	myapi := makeAPI()
	resp, err := jsonapi.HandlerTest(jsonapi.APIHandler(myapi.list).Handler).Get("/api/list", "")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	var ps []Project
	if err := json.Unmarshal(resp.Body.Bytes(), &ps); err != nil {
		t.Fatalf("/api/list does not return valid json format: %s", err)
	}

	if l := len(ps); l != 3 {
		t.Errorf("expected 3 results, got %d", l)
	}

	for _, name := range []string{"a", "b", "c"} {
		p, ok := findProject(ps, name)
		if !ok {
			t.Errorf("expected %s exists, but not found", name)
		}

		expect, _ := findProject(data, name)

		if l, e := len(p.Branches), len(expect.Branches); l != e {
			t.Errorf("expected %d branches in %s, got %d", e, name, l)
		}

		for _, e := range expect.Branches {
			br, ok := findBranch(p, e.Name)
			if !ok {
				t.Errorf("expected branch %s exists in %s, but not found", e.Name, name)
			}
			if br.Owner != e.Owner {
				t.Errorf("expected owner of %s/%s to be %s, got %s", name, e.Name, e.Owner, br.Owner)
			}
			if br.Desc != e.Desc {
				t.Errorf("expected desc of %s/%s to be %s, got %s", name, e.Name, e.Desc, br.Desc)
			}
		}
	}
}
