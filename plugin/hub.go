package plugin

import "github.com/team4yf/yf-fpm-server-go/fpm"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	Namespace string
	// Registered clients.
	Clients map[string]*Client

	// Inbound messages from the clients.
	Send chan *Msg

	// Register requests from the clients.
	Login chan *Client

	// logout requests from clients.
	Logout chan *Client
}

type Msg struct {
	ClientID string
	Payload  []byte
}

func NewHub(ns string) *Hub {
	return &Hub{
		Namespace: ns,
		Send:      make(chan *Msg, 10),
		Login:     make(chan *Client),
		Logout:    make(chan *Client, 10),
		Clients:   make(map[string]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Login:
			h.Clients[client.ID] = client
			fpm.Default().Publish("#ws/connect", client.ID)
		case client := <-h.Logout:
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.Send)
				fpm.Default().Publish("#ws/close", client.ID)
			}
		case message := <-h.Send:
			if message.ClientID == "" {
				// broadcast
				for id, client := range h.Clients {
					select {
					case client.Send <- message.Payload:
					default:
						close(client.Send)
						delete(h.Clients, id)
					}
				}
			} else {
				c, ok := h.Clients[message.ClientID]
				if ok {
					c.Send <- message.Payload
				}
			}
		}
	}
}
