package session

import (
	"context"
	"net/http"
	"time"
)

const SIDKey = "sessid"
const ContextKey = "session"

// ContextData is the data format we store in context
type ContextData struct {
	M    *Manager
	SID  string
	Data *SessionData
}

func (d ContextData) Save() {
	d.M.Save(d.SID, d.Data)
}

// middleware is session middleware
type middleware struct {
	M *Manager
	*http.ServeMux
}

func (m *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	c, err := r.Cookie(SIDKey)
	sid := ""
	var data *SessionData
	if err == nil {
		sid = c.Value
	}

	data = m.M.Get(sid)
	if data == nil {
		sid, data = m.M.Allocate()
	}
	http.SetCookie(w, &http.Cookie{
		Name:     SIDKey,
		Value:    sid,
		Expires:  time.Now().Add(time.Duration(m.M.GetTTL()) * time.Millisecond),
		HttpOnly: true,
	})
	ctx = context.WithValue(ctx, ContextKey, ContextData{
		M:    m.M,
		SID:  sid,
		Data: data,
	})

	m.ServeMux.ServeHTTP(w, r.WithContext(ctx))
}

func Middleware(sess *Manager, mux *http.ServeMux) http.Handler {
	return &middleware{sess, mux}
}
