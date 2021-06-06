package web

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func ServeFile(port int) error {
	// handle font
	box := packr.New("gui", "./gui")
	http.Handle("/", http.FileServer(box))
	http.HandleFunc("/test", updateTest)
	log.Printf("Start server at http://127.0.0.1:%d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	return err
}

func ServeWasm(port int) error {
	// handle font
	box := packr.New("wasm", "./wasm")
	http.Handle("/", http.FileServer(box))
	http.HandleFunc("/test", updateTest)
	log.Printf("Start server at http://127.0.0.1:%d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	return err
}

func updateTest(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		// log.Printf("recv: %s", message)
		links, options, err := parseMessage(message)
		if err != nil {
			log.Println("parseMessage:", err)
			continue
		}
		p := ProfileTest{
			Conn:        c,
			MessageType: mt,
			Links:       links,
			Options:     options,
		}
		go p.testAll(ctx)
		// err = c.WriteMessage(mt, getMsgByte(0, "gotspeed"))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
