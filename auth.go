package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/Patrolavia/jsonapi"
	"github.com/Ronmi/br-tr-backend/session"
	gogitlab "github.com/plouc/go-gitlab-client"
	"golang.org/x/oauth2"
)

type auth struct {
	config  *oauth2.Config
	baseURL string
	path    string
	token   string
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
	token, err := a.config.Exchange(context.TODO(), code)
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
	if sess.Data.User.Id != 0 {
		return result{
			ID:       sess.Data.User.Id,
			Username: sess.Data.User.Username,
			Name:     sess.Data.User.Name,
			Avatar:   sess.Data.User.AvatarUrl,
		}, nil
	}

	// fetch user info
	client := gogitlab.NewGitlab(a.baseURL, a.path, a.token)
	client.Client = a.config.Client(context.TODO(), sess.Data.Token)
	u, err := client.CurrentUser()
	if err != nil {
		return guest, nil
	}

	sess.Data.User = u
	sess.Save()

	return result{
		ID:       u.Id,
		Username: u.Username,
		Name:     u.Name,
		Avatar:   u.AvatarUrl,
	}, nil
}
