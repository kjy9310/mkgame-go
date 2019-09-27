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
		log.Println("original x,y : ", x, y)
		if lastestMove, exist := queMove4Uuid[client.User.Uuid]; exist {
			moveData := lastestMove.Data.(map[string]interface{})
			client.User.Direction = moveData["Direction"].(float64)
			client.User.Speed = moveData["Speed"].(float64)
			moveBeforeAction := float64(lastestMove.Time - lastQueSent)/100
			x+=moveBeforeAction*client.User.Speed*math.Sin(client.User.Direction)
			y+=-1*moveBeforeAction*client.User.Speed*math.Cos(client.User.Direction)
			moveAfterAction := (float64(100) - moveBeforeAction)/100
			x+=moveAfterAction*client.User.Speed*math.Sin(client.User.Direction)
			y+=-1*moveAfterAction*client.User.Speed*math.Cos(client.User.Direction)
			log.Println("xy inside IF",x,y)
		} else {
			log.Println("contiues movement", client.User.Direction)
			x+=client.User.Speed*math.Sin(client.User.Direction)
			y+=-1*client.User.Speed*math.Cos(client.User.Direction)
			log.Println("xy inside IF",x,y)
		}
		newPosition := getPositionFromCoordinate(math.Round(x),math.Round(y))
		log.Println("movement done :", x, y)
		client.User.move(newPosition, client.Map, h)
	}
}

func getCoordinateFromPosition (position int) (float64, float64){
	x:= int(position/10000)
	y:= position-x*10000
	return float64(x), float64(y)
}

func getPositionFromCoordinate (x float64, y float64) int {
	return int(x)*10000+int(y)
}

func (o *Object) move (position int, mapName string, hub *Hub) {
	nX,nY := getCoordinateFromPosition(position)
	if (nX>9999){
		nX=9999
	}else if(nX<1){
		nX=1
	}
	if (nY>9999){
		nY=9999
	}else if (nY<1) {
		nY=1
	}
	newPosition := getPositionFromCoordinate(nX,nY)
	if _, exist := hub.Objects[mapName][newPosition]; exist{
		return
	} else {
		delete(hub.Objects[mapName], o.Position)
		o.Position = newPosition
		hub.Objects[mapName][o.Position] = o
	}
}
