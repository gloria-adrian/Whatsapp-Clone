
package websocket

import(
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r*http.Request) bool {
		return true // Allow all connections by default
	},
}

type client struct {
	conn *websocket.Conn
	swend chan []byte
	
}

type hub struct {
	clients map[*client]bool
Broadcast chan[]byte
Register chan *client
Unregister chan *client 
}

var HubInstance = &hub{
	clients:   make(map[*client]bool),
    Broadcast:     make(chan []byte),
    Register: make(chan *client),
    Unregister: make(chan *client),
}

func (c *client) Read() {
	defer func() {
		HubInstance.Unregister <- 
		c.conn.close()
	}()
	for {
		_, message, err := c.conn.ReadMeassage()
		if err != nil {
			log.Printlmn("read error:", err)
			break
		}
		HubInstance.Broadcast <- message
	}
}

func (c *client) Write() {
	defer func() {
		c. conn.close()
	}()
	for {
		select {
		case message,ok := <-c.Send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register
		h.clients[client] = true
		case cluent := <- h.Unregister
		if _,ok := h.clients[client]; ok {
			delete (h.clients,client)
			close(client.send)
		}
	case message :=<- h.Broadcast:
		for client := range h.clients {
			select {
			case client.Send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
	}

		}
	
	func serveWs(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w,r,nil)
		if err != nil {
			log.Println("upgrade error:",err)
			return
		}
		client := &client{conn: conn,send: make(chan []byte)}
		HubInstance.Register <- client
		go client.Read()
		go client.Write()
	}
