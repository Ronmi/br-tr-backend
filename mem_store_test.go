package main

import (
	"fmt"
	"reflect"
	"testing"
)

func makeMemStore() *MemStore {
	return &MemStore{data}
}

func TestMemStoreList(t *testing.T) {
	store := makeMemStore()

	actual, err := store.List()

	if err != nil {
		t.Fatalf("unexpected error occured when testing MemStore.List: %s", err)
	}

	if !reflect.DeepEqual(actual, data) {
		t.Errorf("incorrect result of MemStore.List, dump it:\n%#v", actual)
	}
}

func TestMemStoreUpdateOwner(t *testing.T) {
	store := makeMemStore()

	if err := store.UpdateOwner("a", "main", "test"); err != nil {
		t.Fatalf("unexpected error occured when testing MemStore.List: %s", err)
	}

	ps, _ := store.List()
	p, ok := findProject(ps, "a")
	if !ok {
		t.Fatal("cannot find modified project after MemStore.UpdateOwner")
	}
	br, ok := findBranch(p, "main")
	if !ok {
		t.Fatal("cannot find modified branch after MemStore.UpdateOwner")
	}

	if br.Owner != "test" {
		t.Errorf("expected owner to be test, got %s", br.Owner)
	}
}

func TestMemStoreUpdateDesc(t *testing.T) {
	store := makeMemStore()

	if err := store.UpdateDesc("a", "main", "test"); err != nil {
		t.Fatalf("unexpected error occured when testing MemStore.List: %s", err)
	}

	ps, _ := store.List()
	p, ok := findProject(ps, "a")
	if !ok {
		t.Fatal("cannot find modified project after MemStore.UpdateDesc")
	}
	br, ok := findBranch(p, "main")
	if !ok {
		t.Fatal("cannot find modified branch after MemStore.UpdateDesc")
	}

	if br.Desc != "test" {
		t.Errorf("expected desc to be test, got %s", br.Desc)
	}
}

func TestMemStoreAddProjects(t *testing.T) {
	store := makeMemStore()
	append := []Project{
		Project{
			Name: "c",
			Branches: []Branch{
				Branch{"main", "ronmi", "stable"},
			},
		},
		Project{
			Name: "d",
			Branches: []Branch{
				Branch{"main", "ronmi", "stable"},
				Branch{"dev", "ronmi", "develop"},
			},
		},
	}

	if err := store.AddProjects(append); err != nil {
		t.Fatal("unexpected error occurs at MemStore.AddProject")
	}

	actual, _ := store.List()

	if l := len(actual); l != 4 {
		t.Errorf("expected we have 4 records, got %d", l)
	}

	test := func(ps []Project, name string, expect Project) (msg string) {
		p, ok := findProject(ps, name)
		if !ok {
			return fmt.Sprintf("expect we have project %s, but not found", name)
		}

		if !reflect.DeepEqual(p, expect) {
			return fmt.Sprintf("content of project %s is unexpected, dumping:\n%#v", name, p)
		}

		return ""
	}

	lst := map[string]Project{
		"a": data[0],
		"b": data[1],
		"c": append[0],
		"d": append[1],
	}
	for n, p := range lst {
		if err := test(actual, n, p); err != "" {
			t.Error(err)
		}
	}
}
