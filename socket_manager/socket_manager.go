package socket_manager

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	mes "terminal-client/message"
)

type SocketManager struct {
	Host             string
	Path             string
	OnNewTextMessage func(message mes.TextMessage)

	conn     *websocket.Conn
	username string
	room     string
}

func (sm *SocketManager) ListenAndRegisterUser(username, room string) error {
	u := url.URL{Scheme: "ws", Host: sm.Host, Path: sm.Host}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
		return err
	}
	defer conn.Close()
	sm.conn = conn

	err = registerUser(conn, username, room)

	if err != nil {
		return err
	}

	sm.username = username
	sm.room = room

	for {
		mt, m, err := conn.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			fmt.Println("error: ", err)
			break
		}

		message := &mes.Message{}
		err = json.Unmarshal(m, message)
		if err != nil {
			fmt.Println("error: ", err)
			break
		}

		sm.handleMessage(message)
	}

	return nil
}

func (sm *SocketManager) SendMessage(text string) {
	sm.conn.WriteJSON(&mes.Message{
		Type:    mes.TypeText,
		Payload: text,
	})
}

func registerUser(conn *websocket.Conn, username, room string) error {
	return conn.WriteJSON(&mes.Message{
		Type: mes.TypeRegister,
		Payload: &mes.RegisterMessage{
			Username: username,
			Room:     room,
		},
	})
}

func (sm *SocketManager) handleMessage(m *mes.Message) {
	switch m.Type {
	case mes.TypeText:
		if textMessageMap, ok := m.Payload.(map[string]interface{}); ok {
			textMessage, err := mes.TextMessageFromMap(textMessageMap)
			if err != nil {
				panic(err)
				fmt.Println("error: ", err)
				return
			}
			sm.OnNewTextMessage(*textMessage)
		} else {
			panic("lskjdf")
			fmt.Println("error: ", "InvalidMessage")
		}
	}
}
