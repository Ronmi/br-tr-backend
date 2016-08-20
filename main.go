package main

import (
	"net/http"

	"github.com/Patrolavia/jsonapi"
)

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
	}}

	myapi := &api{store}

	jsonapi.Register([]jsonapi.API{
		jsonapi.API{"/api/list", myapi.list},
		jsonapi.API{"/api/setOwner", myapi.setOwner},
		jsonapi.API{"/api/setDesc", myapi.setDesc},
	}, nil)

	http.ListenAndServe(":8000", nil)
}
