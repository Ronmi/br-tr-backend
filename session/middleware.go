package session

import (
	"context"
	"net/http"
	"time"
)

const SIDKey = "sessid"

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
	ctx = context.WithValue(ctx, "sessmgr", m.M)
	ctx = context.WithValue(ctx, "sessdata", data)

	m.ServeMux.ServeHTTP(w, r.WithContext(ctx))
}

func Middleware(sess *Manager, mux *http.ServeMux) http.Handler {
	return &middleware{sess, mux}
}
