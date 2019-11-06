package networking

import (
	"../events"
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
	events.GoFuncEvent("ClientRegistry.RunRegistry", h.runRegistry)
}
func (h *ClientRegistry) runRegistry() {
	for {
		select {
		case client := <-h.register:
			events.FuncEvent("ClientRegistry.Register", func() {
				h.clients[client] = true
			})
		case client := <-h.unregister:
			events.FuncEvent("ClientRegistry.Unregister", func() {
				if _, ok := h.clients[client]; ok {
					delete(h.clients, client)
					close(client.send)
				}
			})
		}
	}
}
