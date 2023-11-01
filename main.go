package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Message struct {
    Greeting string `json:"greeting"`
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
    wsConn, err := wsUpgrader.Upgrade(w, r, nil)
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

        fmt.Printf("Message Received: %s\n", msg.Greeting)
    }
}

func main() {
    router := mux.NewRouter()

    router.HandleFunc("/socket", WsEndpoint)

    log.Fatal(http.ListenAndServe(":9100", router))
}