package main

import (
	"github.com/lamnguyencse17/gocketio"
	"github.com/lamnguyencse17/gocketio/events"
	"log"
)

type ISubscribePayload struct {
	VideoId string `json:"video_id"`
	Lang    string `json:"lang"`
}

func main(){
	gocket := gocketio.Gocket{Scheme: "wss", Host: "holodex.net", Path: "api/socket.io/", RawQuery: "EIO=4&transport=websocket"}
	gocket.On("connect", func(data events.CallbackData) {
		log.Println("Connected Successfully")
		subscribePayload := ISubscribePayload{VideoId: "lqhYHycrsHk", Lang: "en"}
		gocket.Emit("subscribe", subscribePayload)
		gocket.On("subscribeSuccess", func (data events.CallbackData){
			log.Println(data.Event)
			log.Println(data.Data)
		})
	})
	quit := make(chan bool)
	go gocket.StartGocket(quit)
	<- quit
}