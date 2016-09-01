package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"golang.org/x/oauth2"

	"github.com/Patrolavia/jsonapi"
	"github.com/Ronmi/br-tr-backend/kvstore"
	"github.com/Ronmi/br-tr-backend/session"
	"github.com/Ronmi/gitlab"
	"github.com/Ronmi/gitlab/webhook"
)

type gitlabConf struct {
	AppURL string `json:"appUrl"`
	URL    string `json:"baseUrl"`
	Path   string `json:"apiPath"`
	Token  string `json:"token"`
	AppID  string `json:"appID"`
	Secret string `json:"appSecret"`
}

func loadGitlabConf(fn string) (ret gitlabConf) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Fatalf("Cannot read config from %s: %s", fn, err)
	}

	if err = json.Unmarshal(data, &ret); err != nil {
		log.Fatalf("Incorrect format of config file %s: %s", fn, err)
	}

	return
}

func main() {
	var addr string
	flag.StringVar(&addr, "http", ":8000", "inet address to bind to")
	flag.Parse()
	store := &MemStore{[]Project{}, &sync.RWMutex{}}
	conf := loadGitlabConf("gitlab.json")
	client := gitlab.FromPAT(conf.URL, conf.Path, conf.Token, nil)
	handler := &webhook.Handler{
		Push:         make(chan webhook.PushEvent),
		MergeRequest: make(chan webhook.MergeRequestEvent),
	}
	http.Handle("/api/webhook", handler)
	(&wh{client, store, handler.Push, handler.MergeRequest}).start()

	myapi := &api{store, client}
	oauth := &oauth2.Config{
		ClientID:     conf.AppID,
		ClientSecret: conf.Secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  conf.URL + "/oauth/authorize",
			TokenURL: conf.URL + "/oauth/token",
		},
		RedirectURL: conf.AppURL + "/api/callback",
		Scopes:      []string{"api"},
	}

	// create session middleware
	mux := http.NewServeMux()
	http.Handle("/api/", session.Middleware(&session.Manager{Store: kvstore.NewMemStore()}, mux))

	// apis
	jsonapi.Register([]jsonapi.API{
		jsonapi.API{"/api/list", myapi.list},
		jsonapi.API{"/api/setOwner", myapi.setOwner},
		jsonapi.API{"/api/setDesc", myapi.setDesc},
		jsonapi.API{"/api/addBranch", myapi.addBranch},
	}, mux)

	// oauth entries
	myauth := &auth{
		config: oauth,
		client: client,
	}
	mux.HandleFunc("/api/callback", myauth.Callback)
	jsonapi.Register([]jsonapi.API{
		jsonapi.API{"/api/auth", myauth.Auth},
		jsonapi.API{"/api/me", myauth.Me},
	}, mux)

	http.ListenAndServe(":8000", nil)
}
