package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/light-bull/lightbull/events"
)

// Based on https://github.com/gorilla/websocket/blob/master/examples/chat/

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var newline = []byte{'\n'}

// WebsocketHandler is the handler for incoming websocket messages
type WebsocketHandler func(ws *WebsocketClient, payload *json.RawMessage)

// WebsocketClient is one client connected via websockets
type WebsocketClient struct {
	eventhub *events.EventHub
	conn     *websocket.Conn
	jwt      *JWTManager

	events chan *events.Event
	send   chan []byte

	handlers map[string]WebsocketHandler

	ID            uuid.UUID // TODO
	Authenticated bool
}

// NewWebsocketClient initalizes a new websocket client and runs the handlers
func NewWebsocketClient(conn *websocket.Conn, eventhub *events.EventHub, jwtmanager *JWTManager) *WebsocketClient {
	client := WebsocketClient{
		eventhub: eventhub,
		conn:     conn,
		jwt:      jwtmanager,

		events: make(chan *events.Event, 256),
		send:   make(chan []byte, 256),

		handlers: make(map[string]WebsocketHandler),

		ID:            uuid.New(), // TODO: make unique
		Authenticated: false,
	}

	go client.readPump()
	go client.writePump()

	client.eventhub.RegisterClient(&client)

	return &client
}

// EventChan is there to implement the EventClient interface
func (client *WebsocketClient) EventChan() chan *events.Event {
	return client.events
}

// AddHandler adds a handler for an incoming message with the specified topic
func (client *WebsocketClient) AddHandler(topic string, handler WebsocketHandler) {
	client.handlers[topic] = handler
}

// SendMessage sends a message over the websocket connection
func (client *WebsocketClient) SendMessage(topic string, payload interface{}) {
	client.events <- events.NewEvent(topic, payload)
}

// SendError sends an error back over the websocket connection
func (client *WebsocketClient) SendError(msg string) {
	type errorFormat struct {
		Msg string `json:"msg"`
	}

	client.events <- events.NewEvent("error", errorFormat{Msg: msg})
}

// Close closes the websocket connection
func (client *WebsocketClient) Close() {
	client.conn.Close()
}

// handleRequest is called after a messages was received over websockets
func (client *WebsocketClient) handleRequest(request []byte) {
	// deserialize topic
	type format struct {
		Topic   string           `json:"topic"`
		Payload *json.RawMessage `json:"payload"`
	}
	data := format{}
	err := json.Unmarshal(request, &data)
	if err != nil {
		client.SendError("Invalid data format")
		return
	}

	// run handler
	handler, ok := client.handlers[data.Topic]
	if ok {
		handler(client, data.Payload)
	} else {
		client.SendError("Unknown message topic")
	}
}

// readPump reads the incoming messages on the websocket connection
func (client *WebsocketClient) readPump() {
	// at end of readPump, close the connection and unregister from evenhub
	defer func() {
		client.eventhub.UnregisterClient(client)
		client.conn.Close()
	}()

	// set max message size
	client.conn.SetReadLimit(maxMessageSize)

	// set read deadline -> end connection if pong is not received within time
	client.updateReadDeadline()

	// and update read deadline everytime a pong was received
	client.conn.SetPongHandler(func(string) error {
		client.updateReadDeadline()
		return nil
	})

	// read messages and give it to eventhub
	for {
		// get the data
		_, data, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
				log.Println("Websocket error: " + err.Error())
			}
			return
		}

		// split messages at newline
		messages := bytes.Split(data, newline)
		for _, message := range messages {
			if len(message) != 0 {
				client.handleRequest(message)
			}
		}
	}
}

// writePump sends the events and other outgoing messages to the client
func (client *WebsocketClient) writePump() {
	ticker := time.NewTicker(pingPeriod)

	// at the end of writePump, close the connection and stop the ticker
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case event, ok := <-client.events:
			// convert event to json bytes for sending
			if !ok {
				// if not, the hub closed the channel, so close the connection
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			data, err := json.Marshal(event)
			if err != nil {
				log.Println("Failed to serialize event for websocket")
			}
			client.send <- data
		case message, ok := <-client.send:
			// we want to send something, so first set the deadline
			client.updateWriteDeadline()

			// check if reading from channel worked
			if !ok {
				return
			}

			// write the message
			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// maybe we have more messages waiting, so send them too
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}

			// close the writer and handle errors
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			// it's time for a ping
			client.updateWriteDeadline()
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}

}

func (client *WebsocketClient) updateReadDeadline() {
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
}

func (client *WebsocketClient) updateWriteDeadline() {
	client.conn.SetWriteDeadline(time.Now().Add(writeWait))
}
