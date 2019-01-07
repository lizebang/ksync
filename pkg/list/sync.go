package list

import (
	"encoding/json"
	"io"
	"os"
	"sync"
)

type (
	Sync struct {
		ps Projects
		tc chan task
		lk *sync.Mutex
		oc *sync.Once
		cl bool
	}
	task func(*Sync)
)

func NewSync() *Sync {
	s := &Sync{
		ps: make(Projects, 0),
		tc: make(chan task, 0),
		lk: &sync.Mutex{},
	}
	go s.deamon()
	return s
}

func (s *Sync) Close() {
	s.oc.Do(func() {
		s.cl = true
		close(s.tc)
	})
}

func (s *Sync) deamon() {
	for t := range s.tc {
		t(s)
	}
}

func (s *Sync) SetProject(name string, proj Project) {
	s.lk.Lock()
	defer s.lk.Unlock()
	if !s.cl {
		s.tc <- task(func(s *Sync) {
			s.ps.Set(name, proj)
		})
	}
}

func (s *Sync) DelProject(name string) {
	s.lk.Lock()
	defer s.lk.Unlock()
	if !s.cl {
		s.tc <- task(func(s *Sync) {
			s.ps.Del(name)
		})
	}
}

func (s *Sync) SetImages(proj, name string, imgs Images) {
	s.lk.Lock()
	defer s.lk.Unlock()
	if !s.cl {
		s.tc <- task(func(s *Sync) {
			if s.ps.IsExist(proj) {
				s.ps[proj].Set(name, imgs)
			}
		})
	}
}

func (s *Sync) DelImages(proj, name string) {
	s.lk.Lock()
	defer s.lk.Unlock()
	if !s.cl {
		s.tc <- task(func(s *Sync) {
			if s.ps.IsExist(proj) {
				s.ps[proj].Del(name)
			}
		})
	}
}

func (s *Sync) SetImage(proj, imgs, name string, img Image) {
	s.lk.Lock()
	defer s.lk.Unlock()
	if !s.cl {
		s.tc <- task(func(s *Sync) {
			if s.ps.IsExist(proj) && s.ps[proj].IsExist(imgs) {
				s.ps[proj][imgs].Set(name, img)
			}
		})
	}
}

func (s *Sync) DelImage(proj, imgs, name string) {
	s.lk.Lock()
	defer s.lk.Unlock()
	if !s.cl {
		s.tc <- task(func(s *Sync) {
			if s.ps.IsExist(proj) && s.ps[proj].IsExist(imgs) {
				s.ps[proj][imgs].Del(name)
			}
		})
	}
}

func (s *Sync) GetImage(proj, imgs, img string) Image {
	s.lk.Lock()
	defer s.lk.Unlock()
	if s.ps.IsExist(proj) && s.ps[proj].IsExist(imgs) && s.ps[proj][imgs].IsExist(img) {
		return s.ps[proj][imgs][img]
	}
	return ""
}

// func (s *Sync) Walk(func()) {
// }

func OpenSync(name string) (*Sync, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := NewSync()
	d := json.NewDecoder(f)
	err = d.Decode(&s.ps)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Sync) WriteTo(w io.Writer) error {
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	return e.Encode(s.ps)
}
