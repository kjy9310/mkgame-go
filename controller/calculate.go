package controller

import (
	"log"
	"math"
)

func (h *Hub) CutObjectsWithRadius(o *Object, mapName string, radius int) map[int]*Object {
	newObjects := make(map[int]*Object)
	x,y :=getCoordinateFromPosition(o.Position)
	xMin := (int(x)-radius)*10000
	xMax := (int(x)+radius)*10000
	yMin := int(y)-radius
	yMax := int(y)+radius
	for position, object := range h.Maps[mapName].Objects {
		if (position > xMin && position < xMax){
			objectY := position%10000
			if (objectY > yMin && objectY < yMax){
				newObjects[position]=object
			}
		}
	}
	return newObjects
}

func (h *Hub) GetObjectWithUuid(uuid string, mapName string) *Object {
	for _, object := range h.Maps[mapName].Objects {
		if object.Uuid == uuid{
			return object
		}
	}
	return nil
}

func (h *Hub) Calculate(sortedQue []wsResponse) {
	queMove4Uuid := map[string]wsResponse{}
	for _, que := range(sortedQue){
		switch que.Action {
		case "move":
			queMove4Uuid[que.Uuid]=que
		case "attack":
			attackData := que.Data.(map[string]interface{})
			client := h.Clients[que.Uuid]
			if client.Cool > 0{
				continue
			}
			attacker := client.User
			log.Println("attacker :",attacker , " attackTarget : ", attackData)
			if (attackData["Target"]==nil){
				continue
			}
			attackTargetUuid := attackData["Target"].(string)
			target := h.GetObjectWithUuid(attackTargetUuid, client.Map)
			if ((int(math.Abs(float64(target.Position/10000 - attacker.Position/10000)))<50) && (int(math.Abs(float64(target.Position%10000 - attacker.Position%10000)))<50)) {
				client.Cool = 10
				target.Hp = target.Hp - attacker.Ap - target.Dp
			}
		}
	}
	lastQueSent := h.LastQueSent
	for _, client := range(h.Clients) {
		x, y := getCoordinateFromPosition(client.User.Position)
		if lastestMove, exist := queMove4Uuid[client.User.Uuid]; exist {
			moveData := lastestMove.Data.(map[string]interface{})
			if (moveData["Direction"]==nil || moveData["Speed"]==nil){
				continue
			}
			client.User.Direction = moveData["Direction"].(float64)
			client.User.Speed = moveData["Speed"].(float64)
			moveBeforeAction := float64(lastestMove.Time - lastQueSent)/100
			x+=moveBeforeAction*client.User.Speed*math.Sin(client.User.Direction)
			y+=-1*moveBeforeAction*client.User.Speed*math.Cos(client.User.Direction)
			moveAfterAction := (float64(100) - moveBeforeAction)/100
			x+=moveAfterAction*client.User.Speed*math.Sin(client.User.Direction)
			y+=-1*moveAfterAction*client.User.Speed*math.Cos(client.User.Direction)
		} else {
			x+=client.User.Speed*math.Sin(client.User.Direction)
			y+=-1*client.User.Speed*math.Cos(client.User.Direction)
		}
		newPosition := getPositionFromCoordinate(math.Round(x),math.Round(y))
		client.User.move(newPosition, client.Map, h)
		client.Cool -= 1
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
	if _, exist := hub.Maps[mapName].Objects[newPosition]; exist{
		return
	} else {
		delete(hub.Maps[mapName].Objects, o.Position)
		o.Position = newPosition
		hub.Maps[mapName].Objects[o.Position] = o
	}
}

func (m *Map) ObjectGenerator(){
	m.Objects[4000400] = &Object{
		Uuid: "tree1",
		Direction: 0,
		Position: 4000400,
		Dp:5,
		Hp:10,
	}
}
