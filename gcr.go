package main

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/google"
)

const (
	gcrprojects = "gcrprojects"
)

func (c *Client) getProjectList(r io.Reader) error {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	bufS := strings.TrimSuffix(string(buf), "\n")
	c.pjs = strings.Split(bufS, "\n")
	return nil
}

func (c *Client) gcrWalk(root string, walkFn google.WalkFunc, options ...google.ListerOption) error {
	repo, err := name.NewRepository(root, name.WeakValidation)
	if err != nil {
		return err
	}

	return google.Walk(repo, walkFn, options...)
}

func (c *Client) gcrWalkFunc(repo name.Repository, tags *google.Tags, err error) error {
	if err != nil {
		return err
	}

	os.MkdirAll(path.Join(c.dir, repo.RepositoryStr()), 0755)

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
