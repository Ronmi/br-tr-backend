package session

import (
	"testing"

	"github.com/Ronmi/br-tr-backend/kvstore"
)

func makeManager() *Manager {
	return &Manager{
		Store: kvstore.NewMemStore(),
		TTL:   500,
	}
}

func TestSessionAllocate(t *testing.T) {
	m := makeManager()

	s1, _ := m.Allocate()
	s2, _ := m.Allocate()

	if s1 == "" || s2 == "" {
		t.Fatal("Manager.Allocate() must not return empty string")
	}

	if l := len(s1); l != SessionIDLength {
		t.Fatalf("Session ID returned from Manager.Allocate() must be %d chars, got %d: %s", SessionIDLength, l, s1)
	}
	if l := len(s2); l != SessionIDLength {
		t.Fatalf("Session ID returned from Manager.Allocate() must be %d chars, got %d: %s", SessionIDLength, l, s2)
	}

	if s1 == s2 {
		t.Fatal("Got same session id from two different calls to Manager.Allocate()")
	}
}

func TestSession(t *testing.T) {
	m := makeManager()

	sid, old := m.Allocate()
	old.Token = "test"
	m.Save(sid, old)

	data := m.Get(sid)
	if data == nil || data.Token != old.Token {
		t.Fatalf("expected fetched data is same as we saved, get %#v", data)
	}

	m.Destroy(sid)
	if data := m.Get(sid); data != nil {
		t.Fatalf("Expected session to be cleared after Manager.Destroy(), get %#v", data)
	}
}
