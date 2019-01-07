package git

import (
	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

type Client struct {
	Name  string
	Email string
	Token string

	Repository *gogit.Repository
}

func (c *Client) PlainClone(dir, url string) (err error) {
	c.Repository, err = gogit.PlainClone(dir, false, &gogit.CloneOptions{
		URL: url,
		Auth: &http.BasicAuth{
			Username: c.Name,
			Password: c.Token,
		},
	})
	return
}

// func (c *Client) Push()  {
// }
