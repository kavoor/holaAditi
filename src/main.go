package main

import (
    "log"
    "net/http"
    "github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) //connected clients
var broadcast = make(chan Message) //broadcastChannel

//Object to updrade HTTP Connections to websockets
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

type Message struct {
    Email       string `json:"email"`
    Username    string `json:"username"`
    Message     string `json:"message"`
}

func main() {
    fs := http.FileServer(http.Dir("../public"))
    http.Handle("/", fs)

    // Configure websocket route
    http.HandleFunc("/ws", handleConnections)

    //start listening for incoming chat messages
    go handleMessages()

    //log
    log.Println("http server started on :8000")
    err := http.ListenAndServe(":8000",nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Fatal(err)
    }

    defer ws.Close()

    //register new client
    clients[ws] = true

    for {
        var msg Message
        err := ws.ReadJSON(&msg)
        if err != nil {
            log.Printf("error: %v" , err)
            delete(clients, ws)
            break
        }
        broadcast <- msg
    }

}
    func handleMessages(){

        for {
            //grab next message in broadcast channel
            msg := <-broadcast
            for client := range clients {
                err := client.WriteJSON(msg)
                if err != nil {
                    log.Printf("error: %v", err)
                    client.Close()
                    delete(clients, client)
                }
            }

        }
    }
