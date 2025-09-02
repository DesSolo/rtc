package etcd

import clientv3 "go.etcd.io/etcd/client/v3"

// OptionFunc ...
type OptionFunc func(c *Provider)

// WithClient ...
func WithClient(client clientv3.Client) OptionFunc {
	return func(c *Provider) {
		c.client = &client
	}
}

// WithPath ...
func WithPath(path string) OptionFunc {
	return func(c *Provider) {
		c.path = path
	}
}
