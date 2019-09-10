package controller

import (
	"log"
)

type Object struct {
	Uuid string
	X int
	Y int
	Direction int
	Speed float32
	Durability int
}

type Map struct {
	Grid [][]Field
	Uuid string
}

type Field struct {
	Type string
}

func (h *Hub) Calculate(sortedQue []wsResponse) {
	queMove4Uuid := map[string]wsResponse{}

	for _, que := range(sortedQue){
		switch que.Action {
		case "move":
			queMove4Uuid[que.Uuid]=que
		}
	}
	lastQueSent := h.LastQueSent
	for client, _ := range(h.Clients) {
		if lastestMove, exist := queMove4Uuid[client.User.Uuid]; exist {
			moveBeforeAction := lastestMove.Time - lastQueSent
			log.Println("movingTime before Action : ", moveBeforeAction)
			moveAfterAction := lastQueSent + 100 - lastestMove.Time
			log.Println("movingTime after action : ", moveAfterAction)
		} else {
			log.Println("cotniues movement", client.User.Direction)
		}
	}
}

