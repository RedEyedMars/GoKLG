package networking

import (
	"Events"
)

type ClientRegistry struct {
	register   chan *Client
	unregister chan *Client
	clients    map[*Client]bool
}

func newRegistry() *ClientRegistry {
	return &ClientRegistry{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *ClientRegistry) run() {
	Events.GoFuncEvent("ClientRegistry.RunRegistry", h.run_registry)
}
func (h *ClientRegistry) run_registry() {
	for {
		select {
		case client := <-h.register:
			Events.FuncEvent("ClientRegistry.Register", func() {
				h.clients[client] = true
			})
		case client := <-h.unregister:
			Events.FuncEvent("ClientRegistry.Unregister", func() {
				if _, ok := h.clients[client]; ok {
					delete(h.clients, client)
					close(client.send)
				}
			})
		}
	}
}
