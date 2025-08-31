package ws

import "github.com/gorilla/websocket"

type Client struct {
	Conn *websocket.Conn
	ChatroomID uint 
	Send chan []byte 
}

type Hub struct {
	Rooms map[uint]map[*Client]bool 
	Register chan *Client 
	Unregister chan *Client 
	Broadcast chan MessageEvent
}

type MessageEvent struct {
	ChatroomID uint 
	Message []byte 
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[uint]map[*Client]bool),
        Broadcast:  make(chan MessageEvent),
        Register:   make(chan *Client),
        Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if h.Rooms[client.ChatroomID] == nil {
				h.Rooms[client.ChatroomID] = make(map[*Client]bool)
			}
			h.Rooms[client.ChatroomID][client] = true 
		case client := <-h.Unregister:
			if _, ok := h.Rooms[client.ChatroomID][client]; ok {
				delete(h.Rooms[client.ChatroomID], client)
				close(client.Send)
			}
		case event := <-h.Broadcast:
			for client := range h.Rooms[event.ChatroomID] {
				select {
				case client.Send <- event.Message:
				default:
					close(client.Send)
					delete(h.Rooms[event.ChatroomID], client)
				}
			}
		}
	}
}