package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/Patrolavia/jsonapi"
	"github.com/Ronmi/br-tr-backend/kvstore"
	"github.com/Ronmi/br-tr-backend/session"
	gogitlab "github.com/plouc/go-gitlab-client"
)

type gitlabConf struct {
	URL   string `json:"baseUrl"`
	Path  string `json:"apiPath"`
	Token string `json:"token"`
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
	store := &MemStore{[]Project{}, &sync.RWMutex{}}
	conf := loadGitlabConf("gitlab.json")

	myapi := &api{store}
	mywebhook := &webhook{&GogitlabProvider{gogitlab.NewGitlab(conf.URL, conf.Path, conf.Token)}, store}

	// create session middleware
	mux := http.NewServeMux()
	http.Handle("/api/", session.Middleware(&session.Manager{Store: kvstore.NewMemStore()}, mux))

	// apis
	jsonapi.Register([]jsonapi.API{
		jsonapi.API{"/api/list", myapi.list},
		jsonapi.API{"/api/setOwner", myapi.setOwner},
		jsonapi.API{"/api/setDesc", myapi.setDesc},
		jsonapi.API{"/api/webhook", mywebhook.entry},
	}, mux)

	http.ListenAndServe(":8000", nil)
}
