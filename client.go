package main

import (
	"os"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/google"

	"github.com/lizebang/ksync/pkg/gcr"
	"github.com/lizebang/ksync/pkg/git"
	"github.com/lizebang/ksync/pkg/list"
	"github.com/lizebang/ksync/pkg/log"
)

type Client struct {
	dir string
	url string
	git *git.Client

	current  list.Projects
	projects list.Projects
}

const (
	projects = "projects.json"
)

func NewClient() *Client {
	return &Client{
		git: &git.Client{},
	}
}

func (c *Client) Init() error {
	c.git.Name = os.Getenv("GIT_NAME")
	if c.git.Name == "" {
		log.Fatal("$GIT_NAME must be set")
	}
	c.git.Email = os.Getenv("GIT_EMAIL")
	if c.git.Email == "" {
		log.Fatal("$GIT_EMAIL must be set")
	}
	c.git.Token = os.Getenv("GIT_TOKEN")
	if c.git.Token == "" {
		log.Fatal("$GIT_TOKEN must be set")
	}
	c.dir = os.Getenv("REPO_DIR")
	if c.dir == "" {
		log.Fatal("$REPO_DIR must be set")
	}
	c.url = os.Getenv("REPO_URL")
	if c.url == "" {
		log.Fatal("$REPO_URL must be set")
	}

	err := c.git.PlainClone(c.dir, c.url)
	if err != nil {
		return err
	}

	c.current = list.NewProjects()
	name := filepath.Join(c.dir, projects)
	if fileExist(name) {
		c.projects, err = list.Open(name)
		if err != nil {
			return err
		}
	} else {
		file, err := os.Create(name)
		if err != nil {
			return err
		}
		file.Close()
		c.projects = list.NewProjects()
	}

	return nil
}

func fileExist(name string) bool {
	_, err := os.Stat(name)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func (c *Client) Run() {
	gcr.Walk("gcr.io/google-containers", c.walk)
}

func (c *Client) walk(repo name.Repository, tags *google.Tags, err error) error {
	return nil
}
