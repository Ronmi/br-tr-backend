package kvstore

// KVStore is atomic key-value store
type KVStore interface {
	// Get value from store, nil if not found or out-dated.
	Get(k string) interface{}
	// Set will insert or update value in store. The value is valid in next ttl miliseconds.
	// Passing 0 or negative number will make the value expire immediately .
	Set(k string, v interface{}, ttl int)
	// SetIf will update a value only if current value matches oldValue.
	SetIf(k string, v, oldValue interface{}, ttl int)
}
