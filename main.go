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
}

var (
    wsUpgrader = websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
    }

    wsConn *websocket.Conn
)

func WsEndpoint(w http.ResponseWriter, r *http.Request) {     
    wsUpgrader.CheckOrigin = func(r *http.Request) bool { 
        return true
     }

    var err error
    wsConn, err = wsUpgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println("could not upgrade: %s\n", err.Error())        
        return
    }

    defer wsConn.Close()

    for {
        var msg Message

        err := wsConn.ReadJSON(&msg)
        if err != nil {
            fmt.Println("error reading Json: %s\n", err.Error())        
            break
        }

        fmt.Printf("Message Received: %s\n", msg.Msg)
        SendMessage(msg.Msg)
    }
}

func SendMessage(msg string) {
    if wsConn == nil {
        fmt.Println("WebSocket connection has not been established")
        return
    }
    
     response := Message{
        Msg: msg,
    }

    if err := wsConn.WriteJSON(response); err != nil {
        fmt.Println("error writing Json: %s\n", err.Error())
        panic(err)
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

    fmt.Printf("Message Received: %s\n", msg.Msg)

    SendMessage(msg.Msg)
}

func main() {
    router := mux.NewRouter()

    router.HandleFunc("/socket", WsEndpoint)
    router.HandleFunc("/send", wsMessage).Methods("POST")


    log.Fatal(http.ListenAndServe(":9100", router))
}