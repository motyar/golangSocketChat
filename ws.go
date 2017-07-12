package main

import (
    "golang.org/x/net/websocket"
    "fmt"
    "log"
    "net/http"
)

var ActiveClients = make(map[Clients]int) // map clients that are connected now

// Client connection consists of the websocket and the client ip
type Clients struct {
    websocket *websocket.Conn
    IP string
}

//Function that handles all the request and creates response
func WShandler(ws *websocket.Conn) {
    var err error

    // cleanout
    defer ws.Close();
    //Here we are creating list of clients that gets connected
    socketClientIP:= ws.Request().RemoteAddr
    socketClient := Clients{ws, socketClientIP} 
    ActiveClients[socketClient] = 0
    log.Println("Total clients live:", len(ActiveClients))
    //For loop to keep it open, It closes after one Send/Recieve
    for {
        var reply string
        var clientMessage string
        if err = websocket.Message.Receive(ws, &reply); err != nil {
            log.Println("Can't receive ",socketClientIP, err.Error())
        }
        fmt.Println("Received back from client: " + reply)

        clientMessage = socketClient.IP + " : " + reply
        //ForEACH client conected, send back the msg to everyone
        for cs, _ := range ActiveClients {
            if err = websocket.Message.Send(cs.websocket, clientMessage); err != nil {
                // It could not send message to a peer
                log.Println("Could not send message to ", cs.IP, err.Error())
            }
        }
    }
}

func main() {
    port:="1212"
    http.Handle("/socket", websocket.Handler(WShandler))
    log.Println("WebSocket started serving at port:"+port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}
