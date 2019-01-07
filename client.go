package main

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/google"
	gogit "gopkg.in/src-d/go-git.v4"

	"github.com/lizebang/ksync/pkg/gcr"
	"github.com/lizebang/ksync/pkg/git"
	"github.com/lizebang/ksync/pkg/log"
)

type Client struct {
	dir string
	url string
	git *git.Client
}

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
	if err != nil && err != gogit.ErrRepositoryAlreadyExists {
		return err
	}
	return nil
}

func (c *Client) Run() {
	repo, err := gogit.PlainOpen(c.dir)
	if err != nil {
		panic(err)
	}
	wt, err := repo.Worktree()
	if err != nil {
		panic(err)
	}
	st, err := wt.Status()
	if err != nil {
		panic(err)
	}
	for key, val := range st {
		println(key, val.Extra, val.Staging, val.Worktree)
	}

	gcr.Walk("gcr.io/google-containers", c.walk)
}

func (c *Client) walk(repo name.Repository, tags *google.Tags, err error) error {
	if err != nil {
		return err
	}

	for _, value := range tags.Children {
		os.MkdirAll(path.Join(c.dir, repo.RepositoryStr(), value), 0755)
	}

	for digest, manifest := range tags.Manifests {
		dockerfile := "FROM "
		file, err := os.OpenFile(filepath.Join(c.dir, repo.RepositoryStr(), strings.TrimPrefix(digest, "sha256:")), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		if len(repo.RegistryStr()) > 0 {
			dockerfile += repo.RegistryStr() + "/"
		}
		dockerfile += repo.RepositoryStr() + "@" + digest
		file.WriteString(dockerfile)
		file.Close()

		for _, tag := range manifest.Tags {
			file, err := os.OpenFile(filepath.Join(c.dir, repo.RepositoryStr(), tag), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				return err
			}
			file.WriteString(dockerfile)
			file.Close()
		}
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
