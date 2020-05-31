package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/techx/playground/config"
	"github.com/techx/playground/db"
	"github.com/techx/playground/socket"
)

func Init() {
	hub := socket.NewHub()

	go hub.Run()
	go db.ListenForUpdates(func(data []byte) {
		var msg map[string]interface{}
		json.Unmarshal(data, &msg)

		switch msg["type"] {
		case "join", "move":
			hub.Send("home", data)
		}
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		socket.ServeWs(hub, w, r)
	})

	r := newRouter()
	http.Handle("/", r)

	config := config.GetConfig()

	fmt.Println("Serving at", config.GetString("server.addr"))
	err := http.ListenAndServe(config.GetString("server.addr"), nil)

	if err != nil {
		panic(err)
	}
}