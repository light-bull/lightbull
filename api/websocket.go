package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/light-bull/lightbull/api/ws"
)

func (api *API) initWebsocket(router *mux.Router) {
	router.HandleFunc("/api/ws", api.handleWebsocketClient)
}

func (api *API) handleWebsocketClient(w http.ResponseWriter, r *http.Request) {
	// authentication is done inside the websocket connection since JS does not support sending the Authentication header here
	enableCors(&w)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	//upgrader.CheckOrigin = func(r *http.Request) bool { return true } // TODO?

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	ws.NewWebsocketClient(conn, api.eventhub, api.jwt)
}
