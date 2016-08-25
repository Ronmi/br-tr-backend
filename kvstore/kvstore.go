package kvstore

import "encoding/json"

// Data is the data fit the store
type Data interface {
	json.Marshaler
	json.Unmarshaler
}

// KVStore is atomic key-value store
type KVStore interface {
	// Get value from store, nil if not found or out-dated.
	Get(k string) Data
	// Set will insert or update value in store. The value is valid in next ttl miliseconds.
	// Passing 0 or negative number will make the value expire immediately .
	Set(k string, v Data, ttl int)
	// SetIf will update a value only if current value matches oldValue. Returns true if successfully updated.
	SetIf(k string, v, oldValue Data, ttl int) bool
}
