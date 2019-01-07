package list

import (
	"encoding/json"
	"io"
	"os"
)

type (
	Projects map[string]Project
	Project  map[string]Images
	Images   map[string]Image
	Image    string
)

func NewProjects() Projects {
	return make(Projects, 0)
}

func (ps Projects) Set(name string, proj Project) {
	ps[name] = proj
}

func (ps Projects) Del(name string) {
	delete(ps, name)
}

func (ps Projects) IsExist(proj string) bool {
	_, ok := ps[proj]
	return ok
}

func NewProject() Project {
	return make(Project, 0)
}

func (p Project) Set(name string, imgs Images) {
	p[name] = imgs
}

func (p Project) Del(name string) {
	delete(p, name)
}

func (p Project) IsExist(imgs string) bool {
	_, ok := p[imgs]
	return ok
}

func NewImages() Images {
	return make(Images, 0)
}

func (is Images) Set(name string, img Image) {
	is[name] = img
}

func (is Images) Del(name string) {
	delete(is, name)
}

func (is Images) IsExist(img string) bool {
	_, ok := is[img]
	return ok
}

func Open(name string) (Projects, error) {
	var ps Projects
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	d := json.NewDecoder(f)
	err = d.Decode(&ps)
	if err != nil {
		return nil, err
	}
	return ps, nil
}

func (ps Projects) WriteTo(w io.Writer) error {
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	return e.Encode(ps)
}
