package events

const SocketioOpen = "0"
const SocketioPing = "2"
const SocketioPong = "3"
const SocketioConnect = "40"
const SocketioEmit = "42"

//const EVENT_SUBSCRIBE = "subscribe"
//const EVENT_SUBSCRIBE_SUCCESS = "subscribeSuccess"

const EventConnect = "connect"

type CallbackData struct {
	Event string
	Data  interface{}
}