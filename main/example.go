package main

import (
	"github.com/lamnguyencse17/gocketio"
	"github.com/lamnguyencse17/gocketio/events"
	"log"
)

func main(){
	gocket := gocketio.Gocket{Scheme: "wss", Host: "holodex.net", Path: "api/socket.io/", RawQuery: "EIO=4&transport=websocket"}
	gocket.On("connect", func(data events.CallbackData) {
		log.Println("Connection Successfully")
		subscribePayload := gocketio.ISubscribePayload{VideoId: "oNOhalk_62w", Lang: "en"}
		gocket.Emit("subscribe", subscribePayload)
		gocket.On("oNOhalk_62w", func (data events.CallbackData){
			log.Println(data.Event)
			log.Println(data.Data)
		})
	})
	quit := make(chan bool)
	go gocket.StartGocket(quit)
	<- quit
}