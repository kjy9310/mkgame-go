package controller

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"time"
	"encoding/json"
	"github.com/google/uuid"
)

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMessageSize = 1024 * 1024
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
}

func ServeWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	c := &Client{
		send: make(chan []byte, maxMessageSize),
		ws: ws,
	}
	id, err := uuid.NewUUID()
	if err !=nil {
		log.Println("UUID creation failed", err)
		return
	}
	c.User.Uuid = id.String()
	// c.User.Speed = 1.0
	c.Map = "start"
	c.User.Position = 5000500
	c.User.Ap = 10
	AHub.Register <- c

	go c.writePump()
	c.readPump()
}


func (c *Client) readPump() {
	defer func() {
		AHub.Unregister <- c
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait));
		return nil
	})

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		var newRequest wsRequest
		err = json.Unmarshal(message, &newRequest)
		if err != nil {
			break
		}
		newRequest.Uuid = c.User.Uuid
		log.Println("request obj : ", newRequest)
		if newRequest.ActionType == "ping" {
			newInput := wsResponse{
				Action : "pong",
				Time : int(time.Now().UnixNano() / int64(time.Millisecond)),
			}
			b, err := json.Marshal(newInput)
			if err != nil {
			        log.Println(err)
				break
			}
			if err := c.write(websocket.TextMessage, b); err != nil {
				break
			}
		} else {
			log.Println(string(message))
			AHub.Receive <- newRequest
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	newInput := wsResponse{
		Action : "login",
		Data : map[string]string{"Uuid":c.User.Uuid},
		Time : int(time.Now().UnixNano() / int64(time.Millisecond)),
	}
	b, err := json.Marshal(newInput)
	if err != nil {
	        log.Println(err)
		return
	}
	if err := c.write(websocket.TextMessage, b); err != nil {
		return
	}
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *Client) write(mt int, message []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, message)
}

