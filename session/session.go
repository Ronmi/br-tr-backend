package session

import (
	cryptoRand "crypto/rand"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/Ronmi/br-tr-backend/kvstore"
)

const SessionIDLength = 16
const DefaultTTL = 3600000 // default 1hr
const chars = "1234567890-qwertyuiopasdfghjklzxcvbnm,.QWERTYUIOPLKJHGFDSAZXCVBNM~#%&_+"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// SessionData defines what kind of data to be stored in session store
type SessionData struct {
	Token string // gitlab access token
}

func (d *SessionData) MarshalJSON() ([]byte, error) {
	return json.Marshal(d)
}
func (d *SessionData) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, d)
}

// Manager manages all session
type Manager struct {
	Store kvstore.KVStore
	TTL   int
}

func (m *Manager) GetTTL() int {
	if m.TTL > 0 {
		return m.TTL
	}
	return DefaultTTL
}

// Get a SessionData
func (m *Manager) Get(sid string) (data *SessionData) {
	ret := m.Store.Get(sid)
	if ret == nil {
		return nil
	}
	return ret.(*SessionData)
}

func (m *Manager) randomString() string {
	ret := make([]byte, SessionIDLength)
	if n, err := cryptoRand.Read(ret); err != nil || n != SessionIDLength {
		// fallback to non-secure random string
		return m.mathRandomString()
	}

	for idx, n := range ret {
		ret[idx] = chars[int(n)%len(chars)]
	}
	return string(ret)
}

func (m *Manager) mathRandomString() string {
	ret := make([]byte, SessionIDLength)
	for idx, _ := range ret {
		ret[idx] = chars[rand.Intn(len(chars))]
	}
	return string(ret)
}

// Allocate a new session
func (m *Manager) Allocate() (sid string, data *SessionData) {
	sid = m.randomString()
	data = &SessionData{}
	ttl := m.GetTTL()

	for !m.Store.SetIf(sid, data, nil, ttl) {
		sid = m.randomString()
	}
	return
}

// Save session data
func (m *Manager) Save(sid string, data *SessionData) {
	m.Store.Set(sid, data, m.GetTTL())
}

// Destroy a session
func (m *Manager) Destroy(sid string) {
	m.Store.Set(sid, nil, 0)
}
