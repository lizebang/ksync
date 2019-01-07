package gcr

import (
	"os"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/google"
)

func List(root string) (*google.Tags, error) {
	repo, err := name.NewRepository(root, name.WeakValidation)
	if err != nil {
		return nil, err
	}

	return google.List(repo)
}

func ListChildren(root string) ([]string, error) {
	repo, err := name.NewRepository(root, name.WeakValidation)
	if err != nil {
		return nil, err
	}

	tags, err := google.List(repo)
	if err != nil {
		return nil, err
	}

	return tags.Children, nil
}

func ListManifests(root string) (map[string]google.ManifestInfo, error) {
	repo, err := name.NewRepository(root, name.WeakValidation)
	if err != nil {
		return nil, err
	}

	tags, err := google.List(repo)
	if err != nil {
		return nil, err
	}

	return tags.Manifests, nil
}

func Walk(root string, walkFn google.WalkFunc, options ...google.ListerOption) error {
	repo, err := name.NewRepository(root, name.WeakValidation)
	if err != nil {
		return err
	}

	return google.Walk(repo, walkFn, options...)
}

func StoreAll(root string) error {
	return Walk(root, store)
}

func store(repo name.Repository, tags *google.Tags, err error) error {
	for _, value := range tags.Children {
		os.MkdirAll(repo.RepositoryStr()+"/"+value, 0755)
	}

	for digest, manifest := range tags.Manifests {
		dockerfile := "FROM "

		file, err := os.OpenFile(repo.RepositoryStr()+"/"+strings.TrimPrefix(digest, "sha256:"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err == nil {
			if len(repo.RegistryStr()) > 0 {
				dockerfile += repo.RegistryStr() + "/"
			}
			dockerfile += repo.RepositoryStr() + "@" + digest
			file.WriteString(dockerfile)
			file.Close()
		}

		for _, tag := range manifest.Tags {
			file, err := os.OpenFile(repo.RepositoryStr()+"/"+tag, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err == nil {
				file.WriteString(dockerfile)
				file.Close()
			}
		}
	}

	return nil
}
