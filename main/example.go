package main

import (
	"encoding/json"
	"github.com/lamnguyencse17/gocketio"
	"log"
)

func main(){
	gocket := gocketio.Gocket{Scheme: "wss", Host: "holodex.net", Path: "api/socket.io/", RawQuery: "EIO=4&transport=websocket"}
	gocket.On("connect", func(empty interface{}) {
		log.Println("Connection Successfully")
	})
	go gocket.StartGocket()
	subscribePayload, err := json.Marshal(gocketio.ISubscribePayload{VideoId: "k_oMkblkB9k", Lang: "en"})
	if err != nil {
		return
	}
	gocket.Emit(gocketio.EVENT_SUBSCRIBE, subscribePayload)
}