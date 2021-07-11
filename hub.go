package main

type Hub struct {
	clients    map[string]*Client
	send       chan Message
	register   chan *Client
	unregister chan *Client
}

type Message struct {
	id   string
	data []byte
}

func newHub() *Hub {
	return &Hub{
		send:       make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.id] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.id]; ok {
				delete(h.clients, client.id)
				close(client.send)
			}
		case message := <-h.send:
			for id, client := range h.clients {
				select {
				case client.send <- message.data:
				default:
					close(client.send)
					delete(h.clients, id)
				}
			}
		}
	}
}
