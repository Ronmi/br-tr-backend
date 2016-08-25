package kvstore

import (
	"strconv"
	"testing"
	"time"
)

type data int

func (d *data) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(*d))), nil
}
func (d *data) UnmarshalJSON(b []byte) error {
	i, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}

	*d = data(i)
	return nil
}

func TestMemKVStore(t *testing.T) {
	s := NewMemStore()

	one := data(1)
	two := data(2)
	three := data(3)
	s.Set("test", &one, 500)
	if v := s.Get("test").(*data); *v != one {
		t.Errorf("expected to get 1, got %#v", *v)
	}
	time.Sleep(500 * time.Millisecond) // wait until expire
	if v := s.Get("test"); v != nil {
		t.Errorf("expected to be expired, got %#v", *(v.(*data)))
	}

	s.Set("test", &one, 500)
	s.SetIf("test", &two, &three, 500)
	if v := s.Get("test").(*data); *v != one {
		t.Errorf("expected to get 1 after failed SetIf, got %#v", *v)
	}

	// wait some time, for test ttl again
	time.Sleep(200 * time.Millisecond) // wait until expire

	s.SetIf("test", &two, &one, 500)
	if v := s.Get("test").(*data); *v != two {
		t.Errorf("expected to get 2 after succeeded SetIf, got %#v", *v)
	}

	// wait until first expire
	time.Sleep(300 * time.Millisecond) // wait until expire
	if v := s.Get("test").(*data); *v != two {
		t.Errorf("expected to get 2 after old expiration time, got %#v", *v)
	}

	// wait until real expire
	time.Sleep(200 * time.Millisecond) // wait until expire
	if v := s.Get("test"); v != nil {
		t.Errorf("expected expire, got %#v", *(v.(*data)))
	}
}

func TestMemKVStoreSetIf(t *testing.T) {
	s := NewMemStore()

	one := data(1)
	two := data(2)

	s.Set("test", &one, 500)
	if s.SetIf("test", &two, &two, 500) {
		t.Fatal("expected SetIf failed, but it success")
	}
	if !s.SetIf("test", &two, &one, 500) {
		t.Fatal("expected SetIf success, but it failed")
	}

}
