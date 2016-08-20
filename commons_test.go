package main

var data = []Project{
	Project{
		Name: "a",
		Branches: []Branch{
			Branch{"main", "ronmi", "stable"},
			Branch{"dev", "ronmi", "develop"},
		},
	},
	Project{
		Name: "b",
		Branches: []Branch{
			Branch{"main", "ronmi", "stable"},
			Branch{"dev", "ronmi", "develop"},
			Branch{"exp", "fraina", "experimental"},
		},
	},
	Project{
		Name: "c",
		Branches: []Branch{
			Branch{"main", "ronmi", "stable"},
			Branch{"dev", "ronmi", "develop"},
		},
	},
}

func findProject(ps []Project, name string) (ret Project, ok bool) {
	for _, p := range ps {
		if p.Name == name {
			return p, true
		}
	}
	return
}
func findBranch(p Project, name string) (ret Branch, ok bool) {
	for _, b := range p.Branches {
		if b.Name == name {
			return b, true
		}
	}
	return
}
