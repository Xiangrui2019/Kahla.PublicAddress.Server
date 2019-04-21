package kahla

import (
	"Kahla.PublicAddress.Server/consts"
	"Kahla.PublicAddress.Server/events"
	"Kahla.PublicAddress.Server/models"
	"encoding/json"
	"github.com/gorilla/websocket"
	"time"
)

type WebSocket struct {
	conn         *websocket.Conn
	Event        chan interface{}
	State        int
	StateChanged chan int
}

func (w *WebSocket) changeState(state int) {
	w.State = state
	select {
	case w.StateChanged <- state:
	default:
	}
}

func NewWebSocket() *WebSocket {
	w := new(WebSocket)
	w.Event = make(chan interface{}, 10)
	w.StateChanged = make(chan int)
	w.changeState(consts.WebSocketStateNew)
	return w
}

func (w *WebSocket) Connect(serverPath string, interrupt chan struct{}) error {
	var err error
	w.conn, _, err = websocket.DefaultDialer.Dial(serverPath, nil)
	if err != nil {
		return err
	}

	defer w.conn.Close()
	w.changeState(consts.WebSocketStateConnected)

	done := make(chan struct{})
	errChan := make(chan error)
	go w.StartReceiveMessage(done, errChan)

	ticker := time.NewTicker(45 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			w.changeState(consts.WebSocketStateDisconnected)
			return nil
		case err := <-errChan:
			w.changeState(consts.WebSocketStateDisconnected)
			return err
		case <-ticker.C:
		case <-interrupt:
			err := w.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				w.changeState(consts.WebSocketStateClosed)
				return err
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			w.changeState(consts.WebSocketStateClosed)
			return nil
		}
	}
}

func (w *WebSocket) StartReceiveMessage(done chan<- struct{}, errChan chan<- error) {
	defer close(done)
	defer close(errChan)
	for {
		_, message, err := w.conn.ReadMessage()
		if err != nil {
			errChan <- err
			return
		}
		event, err := DecodeWebSocketEvent(message)
		if err != nil {
			errChan <- err
			return
		}
		w.Event <- event
	}
}

func DecodeWebSocketEvent(message []byte) (interface{}, error) {
	var err error
	event1 := &events.Event{}
	err = json.Unmarshal(message, event1)
	if err != nil {
		return event1, err
	}
	var event interface{}
	switch event1.Type {
	case models.EventTypeNewMessage:
		event = &events.NewMessageEvent{}
	case models.EventTypeNewFriendRequest:
		event = &events.NewFriendRequestEvent{}
	case models.EventTypeWereDeleted:
		event = &events.WereDeletedEvent{}
	case models.EventTypeFriendAccepted:
		event = &events.FriendAcceptedEvent{}
	case models.EventTypeTimerUpdated:
		event = &events.TimerUpdatedEvent{}
	default:
		return event1, &models.InvalidEventTypeError{event1.Type}
	}
	err = json.Unmarshal(message, event)
	if err != nil {
		return event, err
	}
	return event, nil
}
