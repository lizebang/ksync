package main

import (
	"os"
	"path"
	"path/filepath"
	"time"

	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"

	"github.com/lizebang/ksync/log"
)

type Client struct {
	pjs []string

	dir string
	url string

	name  string
	email string
	token string
	repo  *gogit.Repository
	wktr  *gogit.Worktree
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Init() (err error) {
	c.dir = os.Getenv("REPO_DIR")
	if c.dir == "" {
		log.Fatal("$REPO_DIR must be set")
	}
	c.url = os.Getenv("REPO_URL")
	if c.url == "" {
		log.Fatal("$REPO_URL must be set")
	}

	c.name = os.Getenv("GIT_NAME")
	if c.name == "" {
		log.Fatal("$GIT_NAME must be set")
	}
	c.email = os.Getenv("GIT_EMAIL")
	if c.email == "" {
		log.Fatal("$GIT_EMAIL must be set")
	}
	c.token = os.Getenv("GIT_TOKEN")
	if c.token == "" {
		log.Fatal("$GIT_TOKEN must be set")
	}

	c.repo, err = gogit.PlainClone(c.dir, false, &gogit.CloneOptions{
		URL: c.url,
		Auth: &http.BasicAuth{
			Username: c.name,
			Password: c.token,
		},
	})
	if err != nil && err != gogit.ErrRepositoryAlreadyExists {
		return err
	}
	if err != nil && err == gogit.ErrRepositoryAlreadyExists {
		c.repo, err = gogit.PlainOpen(c.dir)
		if err != nil {
			return err
		}
	}
	c.wktr, err = c.repo.Worktree()
	return
}

func (c *Client) Run() {
	c.run()
}

func (c *Client) run() {
	err := c.initRun()
	if err != nil {
		log.Fatal(err)
	}
	for _, val := range c.pjs {
		err = os.RemoveAll(path.Join(c.dir, val))
		if err != nil {
			log.Fatal(err)
		}
	}

	err = c.gcrRun()
	if err != nil {
		log.Fatal(err)
	}
	err = c.wktr.AddGlob(".")
	if err != nil {
		log.Fatal(err)
	}

	err = c.acrRun()
	if err != nil {
		log.Fatal(err)
	}

	err = c.overRun()
}

func (c *Client) initRun() error {
	err := c.wktr.Pull(&gogit.PullOptions{
		Auth: &http.BasicAuth{
			Username: c.name,
			Password: c.token,
		},
	})
	if err != nil && err != gogit.NoErrAlreadyUpToDate {
		return err
	}

	file, err := os.OpenFile(filepath.Join(c.dir, gcrprojects), os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	return c.getProjectList(file)
}

func (c *Client) gcrRun() error {
	for _, val := range c.pjs {
		err := c.gcrWalk("gcr.io/"+val, c.gcrWalkFunc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) acrRun() error {
	return nil
}

func (c *Client) overRun() error {
	commit, err := c.wktr.Commit("update", &gogit.CommitOptions{
		Author: &object.Signature{
			Name:  c.name,
			Email: c.email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}
	_, err = c.repo.CommitObject(commit)
	if err != nil {
		return err
	}
	return c.repo.Push(&gogit.PushOptions{
		Auth: &http.BasicAuth{
			Username: c.name,
			Password: c.token,
		},
	})
}
