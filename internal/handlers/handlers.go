package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/gorilla/websocket"
	nats "github.com/nats-io/nats.go"
)

var wsChan = make(chan WsPayload)

var clients = make(map[WebSocketConnection]string)

// upgradeConnection is the websocket upgrader from gorilla/websockets
var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WebSocketConnection struct {
	*websocket.Conn
}

type App struct {
	Nc *nats.Conn
	Js nats.JetStreamContext
}

// WsJsonResponse defines the response sent back from websocket
type WsJsonResponse struct {
	Event          string   `json:"event"`
	Message        string   `json:"message"`
	ConnectedUsers []string `json:"connected_users"`
}

// WsPayload defines the websocket request from the client
type WsPayload struct {
	Event    string              `json:"event"`
	Uid      string              `json:"uid"`
	Token    string              `json:"token"`
	Conn     WebSocketConnection `json:"-"`
	RoomName string              `json:"roomName"`
}

// WsEndpoint upgrades connection to websocket
func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected to endpoint")

	var response WsJsonResponse
	response.Message = `<em><small>Connected to server</small></em>`

	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}

	go ListenForWs(&conn)
}

func ListenForWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var payload WsPayload

	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			// do nothing
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

func ListenToWsChannel(app *App) {
	var response WsJsonResponse

	for {
		e := <-wsChan
		clients[e.Conn] = e.Uid
		users := getUserList()
		response.Event = e.Event
		response.ConnectedUsers = users
		broadcastToAll(response)

		byteData, err := json.Marshal(response)
		if err != nil {
			fmt.Println("error in Marshaling", err)
		}

		_, err = app.Js.Publish("ws_sub1", byteData)
		if err != nil {
			fmt.Println("Error in publishing data", err)
		}
	}
}

func getUserList() []string {
	var userList []string
	for _, x := range clients {
		userList = append(userList, x)
	}
	sort.Strings(userList)
	return userList
}

func broadcastToAll(response WsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("websocket err")
			_ = client.Close()
			delete(clients, client)
		}
	}
}
