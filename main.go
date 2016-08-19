package main

import (
	"net/http"

	"github.com/Patrolavia/jsonapi"
)

func main() {
	myapi := &api{map[string]Project{
		"Ronmi/react-toy-router": Project{
			Name: "Ronmi/react-toy-router",
			Branches: []Branch{
				Branch{"main", "ronmi", "stable"},
				Branch{"dev", "ronmi", "develop"},
			},
		},
		"Ronmi/react-promise-visualizer": Project{
			Name: "Ronmi/react-promise-visualizer",
			Branches: []Branch{
				Branch{"main", "ronmi", "stable"},
				Branch{"dev", "ronmi", "develop"},
				Branch{"exp", "fraina", "experimental"},
			},
		},
		"Ronmi/some-go-project": Project{
			Name: "Ronmi/some-go-project",
			Branches: []Branch{
				Branch{"main", "ronmi", "stable"},
				Branch{"dev", "ronmi", "develop"},
			},
		},
	}}

	jsonapi.Register([]jsonapi.API{
		jsonapi.API{"/api/list", myapi.list},
		jsonapi.API{"/api/setOwner", myapi.setOwner},
		jsonapi.API{"/api/setDesc", myapi.setDesc},
	}, nil)

	http.ListenAndServe(":8000", nil)
}
