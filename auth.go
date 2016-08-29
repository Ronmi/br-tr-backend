package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/Patrolavia/jsonapi"
	"github.com/Ronmi/br-tr-backend/session"
	"github.com/Ronmi/gitlab"
	"golang.org/x/oauth2"
)

type auth struct {
	config *oauth2.Config
	client *gitlab.GitLab
}

func (a *auth) Auth(dec *json.Decoder, httpData *jsonapi.HTTP) (interface{}, error) {
	code := strconv.Itoa(rand.Int())
	sess := httpData.Request.Context().Value(session.ContextKey).(session.ContextData)

	sess.Data.ValidateOAUTH = code
	sess.Save()

	return a.config.AuthCodeURL(code, oauth2.AccessTypeOffline), nil
}

func (a *auth) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	state := r.FormValue("state")
	sess := r.Context().Value(session.ContextKey).(session.ContextData)

	if state != sess.Data.ValidateOAUTH || code == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get access token
	token, err := a.config.Exchange(r.Context(), code)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sess.Data.Token = token

	sess.Save()
	http.Redirect(w, r, "/", http.StatusFound)
}

func (a *auth) Me(dec *json.Decoder, httpData *jsonapi.HTTP) (interface{}, error) {
	sess := httpData.Request.Context().Value(session.ContextKey).(session.ContextData)

	type result struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Name     string `json:"name"`
		Avatar   string `json:"avatar"`
	}
	guest := result{
		Username: "guest",
		Name:     "guest",
		Avatar:   "//patrolavia.com/logo64.png",
	}
	if sess.Data.Token == nil {
		return guest, nil
	}
	if sess.Data.User.ID != 0 {
		return result{
			ID:       sess.Data.User.ID,
			Username: sess.Data.User.Username,
			Name:     sess.Data.User.Name,
			Avatar:   sess.Data.User.AvatarURL,
		}, nil
	}

	// fetch user info
	source := a.config.TokenSource(httpData.Request.Context(), sess.Data.Token)
	client := a.client.WithOAuth(source)
	u, err := client.Me()
	if err != nil {
		return guest, nil
	}

	sess.Data.User = u
	sess.Save()

	return result{
		ID:       u.ID,
		Username: u.Username,
		Name:     u.Name,
		Avatar:   u.AvatarURL,
	}, nil
}
