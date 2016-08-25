package kvstore

import (
	"testing"
	"time"
)

func TestMemKVStore(t *testing.T) {
	s := NewMemStore()

	s.Set("test", 1, 500)
	if v := s.Get("test"); v != 1 {
		t.Errorf("expected to get 1, got %#v", v)
	}
	time.Sleep(500 * time.Millisecond) // wait until expire
	if v := s.Get("test"); v != nil {
		t.Errorf("expected to be expired, got %#v", v)
	}

	s.Set("test", 1, 500)
	s.SetIf("test", 2, 3, 500)
	if v := s.Get("test"); v != 1 {
		t.Errorf("expected to get 1 after failed SetIf, got %#v", v)
	}

	// wait some time, for test ttl again
	time.Sleep(200 * time.Millisecond) // wait until expire

	s.SetIf("test", 2, 1, 500)
	if v := s.Get("test"); v != 2 {
		t.Errorf("expected to get 2 after succeeded SetIf, got %#v", v)
	}

	// wait until first expire
	time.Sleep(300 * time.Millisecond) // wait until expire
	if v := s.Get("test"); v != 2 {
		t.Errorf("expected to get 2 after old expiration time, got %#v", v)
	}

	// wait until real expire
	time.Sleep(200 * time.Millisecond) // wait until expire
	if v := s.Get("test"); v != nil {
		t.Errorf("expected expire, got %#v", v)
	}
}

func TestMemKVStoreSetIf(t *testing.T) {
	s := NewMemStore()

	s.Set("test", 1, 500)
	if s.SetIf("test", 2, 2, 500) {
		t.Fatal("expected SetIf failed, but it success")
	}
	if !s.SetIf("test", 2, 1, 500) {
		t.Fatal("expected SetIf success, but it failed")
	}

}
