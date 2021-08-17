package gocketio

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"sync"
)

type Gocket struct {
	Scheme string
	Host     string
	Path       string
	RawQuery   string
	connection *websocket.Conn
	listeners map[string]func(interface{})
	mu sync.Mutex
}

func (gocketObject *Gocket) StartGocket(quit chan bool){

	gocketObject.mu.Lock()
	gocketObject.connection = establishConnection(gocketObject.Scheme, gocketObject.Host, gocketObject.Path, gocketObject.RawQuery)
	gocketObject.mu.Unlock()
	gocketObject.startListening(quit)
}

func (gocketObject *Gocket) On(event string, callBack func(interface{})){
	if gocketObject.listeners == nil {
		gocketObject.listeners = make(map[string]func(interface{}))
	}
	gocketObject.listeners[event] = callBack
}

func (gocketObject *Gocket) Emit(event string, data interface{}){
	payload, err := json.Marshal(data)
	if err != nil {
		return
	}
	subscribeMessage := wrapMessage(SOCKETIO_EMIT, event, string(payload))
	gocketObject.mu.Lock()
	err = gocketObject.connection.WriteMessage(WS_MESSAGE_TYPE, subscribeMessage)
	gocketObject.mu.Unlock()
	if err != nil {
		log.Println(err)
	}
}

func establishConnection(scheme string, host string, path string, rawQuery string) *websocket.Conn {
	u := url.URL{Scheme: scheme, Host: host, Path: path, RawQuery: rawQuery}
	log.Printf("connecting to %s\n", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	return c
}

func (gocketObject *Gocket) startListening(quit chan bool) {
	connection := gocketObject.connection
	defer connection.Close()
	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			quit <- true
		}
		log.Printf("recv: %s", message)
		parsedMessage := parseMessage(message)
		switch parsedMessage.code {
		case SOCKETIO_OPEN:
			{
				connectMessage := []byte(SOCKETIO_CONNECT)
				gocketObject.mu.Lock()
				err := connection.WriteMessage(WS_MESSAGE_TYPE, connectMessage)
				gocketObject.mu.Unlock()
				if err != nil {
					log.Println(err)
					quit <- true
				}
			}
		case SOCKETIO_CONNECT:
			{
				//TODO: Legacy from holovn. Migrating later
				//log.Println("CONNECTION ESTABLISHED")
				//subscribePayload, err := json.Marshal(ISubscribePayload{VideoId: "k_oMkblkB9k", Lang: "en"})
				//if err != nil {
				//	return
				//}
				//subscribeMessage := wrapMessage(SOCKETIO_EMIT, EVENT_SUBSCRIBE, string(subscribePayload))
				//err = connection.WriteMessage(WS_MESSAGE_TYPE, subscribeMessage)

				if callBack, ok := gocketObject.listeners["connect"]; ok {
					// TODO: Generalize this later
					var empty interface{}
					callBack(empty)
				}
			}
		case SOCKETIO_EMIT:
			{
				log.Println("EMIT SUCCESS")
			}
		case SOCKETIO_PING:
			{
				err = pongSocketIO(connection)
				if err != nil {
					log.Println(err)
					quit <- true
				}
			}
		default:
			return
		}
	}
}

func pongSocketIO(connection *websocket.Conn) error {
	pongMessage := []byte(SOCKETIO_PONG)
	err := connection.WriteMessage(WS_MESSAGE_TYPE, pongMessage)
	if err != nil {
		return err
	}
	return nil
}
