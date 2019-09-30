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
	Clients:	 make(map[string]*Client),
	Inputqueues:	make([]wsResponse,0),
	LastQueSent: 0,
	Maps: make(map[string]*Map),
}


var packetPeriod = 100 * time.Millisecond

func (h *Hub) loadMaps(){
	objects := make(map[int]*Object)
	h.Maps = make(map[string]*Map)
	h.Maps["start"] = &Map{Uuid:"start", Objects:objects}
	h.Maps["start"].ObjectGenerator()
}

func (h *Hub) Run() {
	packet := time.NewTicker(packetPeriod)
	ticks := 0
	syncOnTicks := 10
	h.loadMaps()
	for {
		select {
		case c := <-h.Register:
			h.Clients[c.User.Uuid] = c
			h.Maps[c.Map].Objects[c.User.Position]=&c.User
			break

		case c := <-h.Unregister:
			_, ok := h.Clients[c.User.Uuid]
			if _,exist := h.Maps[c.Map].Objects[c.User.Position]; exist {
				delete(h.Maps[c.Map].Objects,c.User.Position)
			}
			if ok {
				delete(h.Clients, c.User.Uuid)
				close(c.send)
			}
			break

		case m := <-h.Receive:
			h.stackReceivedInfo(m)
			break
		case <- packet.C:
			ticks++
			h.Calculate(h.Inputqueues)
			h.LastQueSent = int(time.Now().UnixNano() / int64(time.Millisecond))

			if (len(h.Inputqueues)<1){
				break
			}
			stringPackets, err := json.Marshal(h.Inputqueues)
			if err != nil {
			        log.Println(err)
				break
			}
			h.broadcastMessage(stringPackets)
			h.Inputqueues = []wsResponse{}
			break

		default:
			if (ticks > syncOnTicks){
				ticks = 0
				h.serverStatus = wsResponse{
					Action : "status",
					Time : int(time.Now().UnixNano() / int64(time.Millisecond)),
					Data : map[string]interface{}{"Maps":h.Maps},
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
	for _, c := range h.Clients {
		select {
		case c.send <- message:
			break
		default:
			close(c.send)
			delete(h.Clients, c.User.Uuid)
		}
	}
}
