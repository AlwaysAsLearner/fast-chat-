package ws

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/AlwaysAsLearner/fast-chat/backend/internal/services"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSMessage struct {
	ChatroomID uint   `json:"chatroom_id"`
	UserID     uint   `json:"user_id"`
	Content    string `json:"content"`
}

func WSServe(hub *Hub, msgService *services.MessageService, w http.ResponseWriter, r *http.Request) {
	chatroomIDParam := r.URL.Query().Get("chatroom_id") // get the chatroom id
	chatroomID, _ := strconv.ParseUint(chatroomIDParam, 10, 64)

	// upgrade to websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// create websocket client (to join/send/recieve message in a chatroom through hub)
	client := &Client{
		Conn:       conn,
		ChatroomID: uint(chatroomID),
		Send:       make(chan []byte, 256), // channel with default size of 256 bit
	}

	// add the client to our hub
	hub.Register <- client

	// sending message
	go WritePump(client)

	// receiving message
	go ReadPump(client, hub, msgService)
}

// write the message owner
func WritePump(client *Client) {
	defer client.Conn.Close()
	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			client.Conn.WriteMessage(websocket.TextMessage, message)
		}

	}
}

// read from the
func ReadPump(client *Client, hub *Hub, msgService *services.MessageService) {
	// make sure to unregister client before closing the pump
	defer func() {
		hub.Unregister <- client
		client.Conn.Close()
	}()

	// loop to read message
	for {
		// read message from the connection
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			break
		}

		// unmarshal the jsonified message
		var incomingMsg WSMessage
		if err := json.Unmarshal(message, &incomingMsg); err != nil {
			continue
		}

		saved, _ := msgService.SendMessage(incomingMsg.Content, incomingMsg.UserID, incomingMsg.ChatroomID)
		outgoingMsg, _ := json.Marshal(saved)
		hub.Broadcast <- MessageEvent{
			ChatroomID: incomingMsg.ChatroomID,
			Message:    outgoingMsg,
		}

	}

}
