package controller

import (
	"log"
	"math"
)

type Object struct {
	Uuid string
// position : 4 digit each  xxxxyyyy 8digit int
	Position int
	Direction float64
	Speed float64
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
	log.Println("Calculate : ", sortedQue)
	for _, que := range(sortedQue){
		switch que.Action {
		case "move":
			queMove4Uuid[que.Uuid]=que
		}
	}
	lastQueSent := h.LastQueSent
	for client, _ := range(h.Clients) {
		x, y := getCoordinateFromPosition(client.User.Position)
		if lastestMove, exist := queMove4Uuid[client.User.Uuid]; exist {
			moveData := lastestMove.Data.(map[string]interface{})
			newDirection := moveData["Direction"].(float64)
			client.User.Speed = moveData["Speed"].(float64)
			moveBeforeAction := lastestMove.Time - lastQueSent
			x+=float64(moveBeforeAction)*client.User.Speed*math.Cos(client.User.Direction)
			y+=float64(moveBeforeAction)*client.User.Speed*math.Sin(client.User.Direction)
			moveAfterAction := float64(lastQueSent + 100 - lastestMove.Time) * -1.0
			x+=moveAfterAction*client.User.Speed*math.Cos(newDirection)
			y+=moveAfterAction*client.User.Speed*math.Sin(newDirection)
		} else {
			log.Println("cotniues movement", client.User.Direction)
			x+=float64(100)*client.User.Speed*math.Cos(client.User.Direction)
			y+=float64(100)*client.User.Speed*math.Sin(client.User.Direction)
		}
		newPosition := getPositionFromCoordinate(x,y)
		client.User.move(newPosition)
	}
}

func getCoordinateFromPosition (position int) (float64, float64){
	x:= int(position/10000)
	y:= position-x
	return float64(x), float64(y)
}

func getPositionFromCoordinate (x float64, y float64) int {
	return int(x)*10000+int(y)
}

func (o *Object) move (position int) {
	x:= int(o.Position/10000)
	y:= o.Position-x
	log.Println("original : ", x, y)
	o.Position = position
	x= int(position/10000)
	y= position-x
	log.Println("new position : ", x, y)
}
