package controller

import(
	"sync"
	"github.com/gorilla/websocket"
)

type Hub struct {
	Clients map[string]*Client
	Receive chan wsRequest
	Register chan *Client
	Unregister chan *Client
	Inputqueues []wsResponse
	LastQueSent int
	serverStatus wsResponse
	Maps	map[string]*Map
}

type Client struct {
	ws *websocket.Conn
	send chan []byte
	mu	sync.Mutex
	Map	string
	User Object
	Cool	int
}

type wsResponse struct {
	Time int
	Action string
	Data interface{}
	Uuid string
}

type wsRequest struct {
	ActionType string		`json:"actionType"`
	Value interface{}		`json:"value"`
	Time int			`json:"time"`
	Uuid string
}

type Object struct {
	Uuid string
// position : 4 digit each  xxxxyyyy 8digit int	
	Position int
	Direction float64
	Speed float64
	Ap	int
	Dp	int
	Hp	int
}

type Map struct {
	Uuid string
	Objects map[int]*Object
}


