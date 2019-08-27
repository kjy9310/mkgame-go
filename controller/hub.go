package controller

import (
//	"time"
//	"encoding/json"
//	"log"
//	"strconv"
)

var AHub = Hub{
	Receive:   make(chan string),
	Register:    make(chan *Client),
	Unregister:  make(chan *Client),
	Clients:	 make(map[*Client]bool),
	Inputqueues:	make(map[string][]QueDatum),
	content:	 "",
}

type Hub struct {
	Clients map[*Client]bool
	Receive chan string
	Register chan *Client
	Unregister chan *Client
	Inputqueues map[string][]QueDatum
	content string
}

type QueDatum struct {
	Time int
	Action string
}

type wsRequest struct {
	ActionType string	`json:"actionType"`
	Value string		`json:"value"`
	Time int		`json:"time"`
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

		case m := <-h.Receive:
			h.content = string(m)
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
