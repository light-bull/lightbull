package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/light-bull/lightbull/api/utils"
	"github.com/light-bull/lightbull/events"
)

func (api *API) initWS(router *mux.Router) {
	router.HandleFunc("/api/ws", api.handleWebsocketClient)
}

func (api *API) handleWebsocketClient(w http.ResponseWriter, r *http.Request) {
	// authentication is done inside the websocket connection since JS does not support sending the Authentication header here
	utils.EnableCors(&w)

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

	client := utils.NewWebsocketClient(conn, api.eventhub)
	client.AddHandler("identify", api.handleWSIdentify)
	client.AddHandler("parameter", api.handleWSParameter)
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

func (api *API) handleWSParameter(ws *utils.WebsocketClient, payload *json.RawMessage) {
	if !ws.Authenticated() {
		ws.SendError("Unauthenticated")
		return
	}

	// get parameter id and new value
	type payloadFormat struct {
		ID    string           `json:"id"`
		Value *json.RawMessage `json:"value"`
	}
	data := payloadFormat{}
	err := json.Unmarshal(*payload, &data)
	if err != nil {
		ws.SendError("Invalid data format")
	}

	// get parameter
	_, _, _, parameter := api.shows.FindParameter(data.ID)
	if parameter == nil {
		ws.SendError("Invalud or unknown ID")
		return
	}

	// update value
	err = parameter.SetFromJSON(*data.Value)
	if err != nil {
		ws.SendError("Failed to set parameter: " + err.Error())
		return
	}

	// trigger event
	api.eventhub.PublishNew(events.ParameterChanged, &parameter, nil, uuid.Nil)
}
