package main

type clientRegistry struct {
	clients map[string]*remoteSite
}

func newRegistry() *clientRegistry {
	return &clientRegistry{
		clients: map[string]*remoteSite{},
	}
}

func (r *clientRegistry) getClient(id string) *remoteSite {
	return r.clients[id]
}

func (r *clientRegistry) addClient(c *remoteSite) {
	c.registry = r
	r.clients[c.id] = c
}

func (r *clientRegistry) removeClient(c *remoteSite) {
	delete(r.clients, c.id)
}
