package ws

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/light-bull/lightbull/events"
)

func (client *WebsocketClient) handleRequest(request []byte) {
	// deserialize topic
	type format struct {
		Topic   string           `json:"topic"`
		Payload *json.RawMessage `json:"payload"`
	}
	data := format{}
	err := json.Unmarshal(request, &data)
	if err != nil {
		client.sendError("Invalid data format")
	}

	// handle message
	switch data.Topic {
	case "identify":
		// get token
		type payloadFormat struct {
			Token string `json:"token"`
		}
		payload := payloadFormat{}
		err := json.Unmarshal(*data.Payload, &payload)
		if err != nil {
			client.sendError("Invalid data format")
		}

		// check token
		// TODO (token can be "")
		//client.events <- events.NewEvent("unidentified", nil)

		// return client id
		type responseFormat struct {
			ID uuid.UUID `json:"connectionId"`
		}
		client.events <- events.NewEvent("identified", responseFormat{ID: client.id})
	default:
		client.sendError("Unknown message topic")
	}
}

func (client *WebsocketClient) sendError(msg string) {
	type errorFormat struct {
		Msg string `json:"msg"`
	}

	client.events <- events.NewEvent("error", errorFormat{Msg: msg})
}
