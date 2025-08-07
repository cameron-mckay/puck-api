package websocket

import "sync"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan []byte

	// register requests from the clients.
	register chan *Client

	// unregister requests from clients.
	unregister chan *Client
	close      chan struct{}
	once       sync.Once
}

var MessageHub *Hub

func Init() {
	MessageHub = &Hub{
		Broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		close:      make(chan struct{}),
		once:       sync.Once{},
	}
	go MessageHub.run()
}

func Close() {
	MessageHub.once.Do((func() {
		close(MessageHub.close)
	}))
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.Broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		case <-h.close:
			for client := range h.clients {
				close(client.send)
				delete(h.clients, client)
			}
			return
		}
	}
}
