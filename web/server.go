package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/xxf098/lite-proxy/request"
)

var upgrader = websocket.Upgrader{}

func ServeFile() {
	http.Handle("/", http.FileServer(http.Dir("web/gui/")))
	http.HandleFunc("/test", updateTest)
	fmt.Println("Start server at http://127.0.0.1:10871")
	http.ListenAndServe(":10871", nil)
}

func updateTest(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, getMsgByte(0, "started"))
		err = c.WriteMessage(mt, getMsgByte(0, "gotserver"))
		elapse, err := request.PingLink(string(message))
		err = c.WriteMessage(mt, getMsgByte(0, "gotping", elapse))
		err = c.WriteMessage(mt, getMsgByte(0, "gotspeed"))
		err = c.WriteMessage(mt, getMsgByte(0, "eof"))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

type Message struct {
	ID       int    `json:"id"`
	Info     string `json:"info"`
	Remarks  string `json:"remarks"`
	Group    string `json:"group"`
	Ping     int64  `json:"ping"`
	Speed    string `json:"speed"`
	MaxSpeed string `json:"maxspeed"`
}

func getMsgByte(id int, typ string, option ...interface{}) []byte {
	msg := Message{ID: id, Info: typ}
	switch typ {
	case "gotserver":
		msg.Remarks = "Server 1"
		msg.Group = "Group 1"
	case "gotping":
		msg.Remarks = "Server 1"
		msg.Group = "Group 1"
		var ping int64
		if len(option) > 0 {
			if v, ok := option[0].(int64); ok {
				ping = v
			}
		}
		msg.Ping = ping
	case "gotspeed":
		msg.Remarks = "Server 1"
		msg.Group = "Group 1"
		msg.Speed = "100.00B"
		msg.MaxSpeed = "100.00B"
	}
	b, _ := json.Marshal(msg)
	return b
}