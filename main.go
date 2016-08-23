package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/Patrolavia/jsonapi"
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
	store := &MemStore{[]Project{
		Project{
			Name: "Ronmi/react-toy-router",
			Branches: []Branch{
				Branch{"main", "ronmi", "stable"},
				Branch{"dev", "ronmi", "develop"},
			},
		},
		Project{
			Name: "Ronmi/react-promise-visualizer",
			Branches: []Branch{
				Branch{"main", "ronmi", "stable"},
				Branch{"dev", "ronmi", "develop"},
				Branch{"exp", "fraina", "experimental"},
			},
		},
		Project{
			Name: "Ronmi/some-go-project",
			Branches: []Branch{
				Branch{"main", "ronmi", "stable"},
				Branch{"dev", "ronmi", "develop"},
			},
		},
	}, &sync.RWMutex{}}
	glconf := loadGitlabConf("gitlab.json")

	myapi := &api{store}
	mywebhook := &webhook{gogitlab.NewGitlab(glconf.URL, glconf.Path, glconf.Token), store}

	jsonapi.Register([]jsonapi.API{
		jsonapi.API{"/api/list", myapi.list},
		jsonapi.API{"/api/setOwner", myapi.setOwner},
		jsonapi.API{"/api/setDesc", myapi.setDesc},
		jsonapi.API{"/api/webhook", mywebhook.entry},
	}, nil)

	http.ListenAndServe(":8000", nil)
}
