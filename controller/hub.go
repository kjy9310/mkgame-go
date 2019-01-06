package controller

var ChatHub = Hub{
	Broadcast:   make(chan string),
	Register:    make(chan *Client),
	Unregister:  make(chan *Client),
	Clients: 	 make(map[*Client]bool),
	content:  	 "",
}

type Hub struct {
	Clients map[*Client]bool
	Broadcast chan string
	Register chan *Client
	Unregister chan *Client

	content string
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Register:
			h.Clients[c] = true
			c.send <- []byte(h.content)
			break

		case c := <-h.Unregister:
			_, ok := h.Clients[c]
			if ok {
				delete(h.Clients, c)
				close(c.send)
			}
			break

		case m := <-h.Broadcast:
			h.content = m
			h.broadcastMessage()
			break
		}
	}
}

func (h *Hub) broadcastMessage() {
	for c := range h.Clients {
		select {
		case c.send <- []byte(h.content):
			break

		// We can't reach the Client
		default:
			close(c.send)
			delete(h.Clients, c)
		}
	}
}
