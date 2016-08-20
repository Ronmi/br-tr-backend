package main

import (
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
