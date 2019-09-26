package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/light-bull/lightbull/api/utils"
)

func (api *API) initWS(router *mux.Router) {
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

	client := utils.NewWebsocketClient(conn, api.eventhub, api.jwt)
	client.AddHandler("identify", api.handleWSIdentify)
}

func (api *API) handleWSIdentify(ws *utils.WebsocketClient, payload *json.RawMessage) {
	// get token
	type payloadFormat struct {
		Token string `json:"token"`
	}
	data := payloadFormat{}
	err := json.Unmarshal(*payload, &data)
	if err != nil {
		ws.SendError("Invalid data format")
	}

	// check token
	if api.jwt.Check(data.Token) {
		// return client id
		type responseFormat struct {
			ID uuid.UUID `json:"connectionId"`
		}
		ws.SendMessage("identified", responseFormat{ID: ws.ID()})
		ws.SetAuthenticated(true)
	} else {
		ws.SendMessage("unidentified", nil)
		//ws.Close() // TODO
	}
}
