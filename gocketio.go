package gocketio

import (
	"github.com/gorilla/websocket"
	"github.com/lamnguyencse17/gocketio/events"
	"github.com/lamnguyencse17/gocketio/message"
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
	listeners map[string]func(data events.CallbackData)
	mu sync.Mutex
}

func (gocketObject *Gocket) StartGocket(quit chan bool){

	gocketObject.mu.Lock()
	gocketObject.connection = establishConnection(gocketObject.Scheme, gocketObject.Host, gocketObject.Path, gocketObject.RawQuery)
	gocketObject.mu.Unlock()
	gocketObject.startListening(quit)
}

func (gocketObject *Gocket) On(event string, callBack func(data events.CallbackData)){
	if gocketObject.listeners == nil {
		gocketObject.listeners = make(map[string]func(data events.CallbackData))
	}
	gocketObject.listeners[event] = callBack
}

func (gocketObject *Gocket) Emit(event string, data interface{}){
	subscribeMessage := message.WrapMessage(events.SocketioEmit, event, data)
	gocketObject.mu.Lock()
	err := gocketObject.connection.WriteMessage(message.WsMessageType, subscribeMessage)
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
		_, rawMessage, err := connection.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			quit <- true
		}
		log.Printf("recv: %s", rawMessage)
		parsedMessage := message.ParseMessage(rawMessage)
		switch parsedMessage.Code {
		case events.SocketioOpen:
			{
				connectMessage := []byte(events.SocketioConnect)
				gocketObject.mu.Lock()
				err := connection.WriteMessage(message.WsMessageType, connectMessage)
				gocketObject.mu.Unlock()
				if err != nil {
					log.Println(err)
					quit <- true
				}
			}
		case events.SocketioConnect:
			{
				//TODO: Legacy from holovn. Migrating later
				//log.Println("CONNECTION ESTABLISHED")
				//subscribePayload, err := json.Marshal(ISubscribePayload{VideoId: "k_oMkblkB9k", Lang: "en"})
				//if err != nil {
				//	return
				//}
				//subscribeMessage := wrapMessage(SOCKETIO_EMIT, EVENT_SUBSCRIBE, string(subscribePayload))
				//err = connection.WriteMessage(WS_MESSAGE_TYPE, subscribeMessage
				if callBack, ok := gocketObject.listeners[events.EventConnect]; ok {
					callBack(events.CallbackData{Event: events.EventConnect})
				}
			}
		case events.SocketioEmit:
			{
				unwrappedMessage := message.UnwrapMessage(parsedMessage.Data)
				if callBack, ok := gocketObject.listeners[unwrappedMessage.Event]; ok {
					callBack(events.CallbackData{Event: unwrappedMessage.Event, Data: unwrappedMessage.Data})
				}
			}
		case events.SocketioPing:
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
	pongMessage := []byte(events.SocketioPong)
	err := connection.WriteMessage(message.WsMessageType, pongMessage)
	if err != nil {
		return err
	}
	return nil
}
