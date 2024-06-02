package main

import (
	"log"
	"net/http"
	"whatsapp-clone/pkg/handlers"
	"whatsapp-clone/pkg/websocket"
)

func main() {
	go websocket.HubInstance.Run()
	http.HandleFunc("/", handlers.ServeHome)
	http.HandleFunc("/ws", handlers.ServeWs)
	log.Fatal(http.ListenAndServe(":8081", nil))

}
