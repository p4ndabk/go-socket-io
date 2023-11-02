package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Message struct {
    Msg string `json:"msg"`
    Name string `json:"name"`
}

var (
    wsUpgrader = websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
    }

    wsConns = make(map[string]map[*websocket.Conn]bool)
)

func WsEndpoint(w http.ResponseWriter, r *http.Request) {     
    wsUpgrader.CheckOrigin = func(r *http.Request) bool { 
        return true
     }

    var err error
    wsConn, err := wsUpgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println("could not upgrade: %s\n", err.Error())        
        return
    }

    vars := mux.Vars(r)
	uuid := vars["uuid"]
    if wsConns[uuid] == nil {
		wsConns[uuid] = make(map[*websocket.Conn]bool)
	}

    wsConns[uuid][wsConn] = true

    defer wsConn.Close()

    for {
        var msg Message

        err := wsConn.ReadJSON(&msg)
        if err != nil {
            fmt.Println("error reading Json: %s\n", err.Error())        
            break
        }

        fmt.Printf("Message Received: %s\n", msg.Msg)
        SendMessage(uuid, msg.Msg, msg.Name)
    }
}

func SendMessage(uuid string, msg string, name string) {
    if wsConns[uuid] == nil {
        fmt.Println("WebSocket connection has not been established")
        return
    }
    
     response := Message{
        Msg: msg,
        Name: name,
    }
    for wsConn := range wsConns[uuid] {
        err := wsConn.WriteJSON(response)
        if err != nil {
            fmt.Println("error writing JSON:", err)
            panic(err)
        }
    } 

    fmt.Printf("Message Sent: %s\n", response.Msg)
}

func wsMessage(w http.ResponseWriter, r *http.Request) {
    var msg Message

    err := json.NewDecoder(r.Body).Decode(&msg)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	vars := mux.Vars(r)
	uuid := vars["uuid"]

    fmt.Printf("Message Received: %s\n", msg.Msg)

    SendMessage(uuid, msg.Msg, msg.Name)
}

func main() {
    router := mux.NewRouter()

    router.HandleFunc("/socket/{uuid}", WsEndpoint)
    router.HandleFunc("/send/{uuid}", wsMessage).Methods("POST")


    log.Fatal(http.ListenAndServe(":9100", router))
}