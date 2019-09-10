package controller

import (
	"encoding/json"
	"log"
//	"strconv"
	"time"
	"sort"
)

var AHub = Hub{
	Receive:   make(chan wsRequest),
	Register:    make(chan *Client),
	Unregister:  make(chan *Client),
	Clients:	 make(map[*Client]bool),
	Inputqueues:	make([]wsResponse,0),
	LastQueSent: 0,
	Maps: make(map[string]Map),
	Objects : make(map[string]map[string]Object),
}


var packetPeriod = 100 * time.Millisecond

type Hub struct {
	Clients map[*Client]bool
	Receive chan wsRequest
	Register chan *Client
	Unregister chan *Client
	Inputqueues []wsResponse
	LastQueSent int
	serverStatus wsResponse
	Maps	map[string]Map
	Objects map[string]map[string]Object
}

type wsResponse struct {
	Time int
	Action string
	Data interface{}
	Uuid string
}

type wsRequest struct {
	ActionType string	`json:"actionType"`
	Value string		`json:"value"`
	Time int		`json:"time"`
	Uuid string
}

func (h *Hub) Run() {
	packet := time.NewTicker(packetPeriod)
	ticks := 0
	syncOnTicks := 10
	for {
		select {
		case c := <-h.Register:
			h.Clients[c] = true
			break

		case c := <-h.Unregister:
			_, ok := h.Clients[c]
			if ok {
				delete(h.Clients, c)
				close(c.send)
			}
			break

		case m := <-h.Receive:
			h.stackReceivedInfo(m)
			break
		case <- packet.C:
			ticks++
			if (len(h.Inputqueues)<1){
				break
			}
			h.Calculate(h.Inputqueues)
			stringPackets, err := json.Marshal(h.Inputqueues)
			if err != nil {
			        log.Println(err)
				break
			}
			h.broadcastMessage(stringPackets)
			h.LastQueSent = h.Inputqueues[len(h.Inputqueues)-1].Time
			h.Inputqueues = []wsResponse{}
			break

		default:
			if (ticks > syncOnTicks){
				ticks = 0
				h.serverStatus = wsResponse{
					Action : "status",
					Time : int(time.Now().UnixNano() / int64(time.Millisecond)),
					Data : map[string]interface{}{"test":Object{}},
				}
				stringStatus, err := json.Marshal(h.serverStatus)
				if err != nil {
				        log.Println(err)
					break
				}
				h.broadcastMessage(stringStatus)
			}
		}
	}
}

func (h *Hub) stackReceivedInfo(singleRequest wsRequest){
	if (singleRequest.Time < h.LastQueSent){
		return
	}
	singleQue := wsResponse{
		Action : singleRequest.ActionType,
		Time : singleRequest.Time,
		Data : singleRequest.Value,
		Uuid : singleRequest.Uuid,
	}
	h.Inputqueues = append(h.Inputqueues, singleQue)
	sort.Slice(h.Inputqueues, func(i, j int) bool {
		return h.Inputqueues[i].Time < h.Inputqueues[j].Time
	})
}

func (h *Hub) broadcastMessage(message []byte) {
	for c := range h.Clients {
		select {
		case c.send <- message:
			break
		default:
			close(c.send)
			delete(h.Clients, c)
		}
	}
}
