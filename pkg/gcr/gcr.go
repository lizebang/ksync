package gcr

import (
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
