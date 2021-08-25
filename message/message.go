package message

import (
	"encoding/json"
	"log"
	"strconv"
)

type IParsedMessage struct {
	Code string
	Data string
}

type IUnwrappedMessage struct {
	Event string
	Data interface{}
}

const WsMessageType = 1

func ParseMessage(message []byte) (parsedMessage IParsedMessage) {
	messageString := string(message)
	if len(messageString) == 1 {
		return IParsedMessage{Code: string(messageString)}
	}
	_, err := strconv.Atoi(string(messageString[1]))
	if err != nil {
		return IParsedMessage{Code: string(messageString[0]), Data: messageString[1:]}
	}
	return IParsedMessage{Code: messageString[:2], Data: messageString[2:]}
}

func WrapMessage(socketIOEvent string, event string, data interface{}) (wrappedMessage []byte) {
	var messageBody [2]interface{}
	messageBody[0] = event
	messageBody[1] = data
	bodyText, err := json.Marshal(messageBody)
	if err != nil {
		return nil
	}
	return append([]byte(socketIOEvent),bodyText...)
}

func UnwrapMessage(message string)(unwrappedMessage IUnwrappedMessage){
	var initialMessage = make([]interface{}, 0)
	err := json.Unmarshal([]byte(message), &initialMessage)
	if err != nil {
		log.Fatalln(err)
	}
	unwrappedMessage.Event, _ = initialMessage[0].(string)
	unwrappedMessage.Data, _ = initialMessage[1].(map[string]interface{})
	return unwrappedMessage
}