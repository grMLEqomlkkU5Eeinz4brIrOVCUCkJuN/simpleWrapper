package core

import (
	users "github.com/simpleWrapper/core/resources"

	"github.com/jaxron/axonet/pkg/client"
)

const defaultBaseURL = "http://localhost:3000/api/v1"

type core struct {
	client  *client.Client
	baseURL string
	users   *users.Users
}

func New() *core {
	c := client.NewClient()
	return &core{
		client:  c,
		baseURL: defaultBaseURL,
	}
}

func (c *core) SetBaseURL(url string) {
	c.baseURL = url
	c.users = nil
}

func (c *core) GetBaseURL() string {
	return c.baseURL
}

func (c *core) GetClient() *client.Client {
	return c.client
}

func (c *core) Users() *users.Users {
	if c.users == nil {
		c.users = users.NewUsers(c.client, c.baseURL)
	}
	return c.users
}
