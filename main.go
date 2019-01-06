package main

import (
	"mkgame-go/web"
	"log"
)

func main() {
	log.Println("start")
	web.ServerOn()
	log.Println("done")
}
