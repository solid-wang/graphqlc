package graphqlc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/url"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type Subscription interface {
	Run(ctx context.Context)
	Stop()
	ResultChan() <-chan Decoder
}

type Subscribe struct {
	url     *url.URL
	header  map[string]string
	payload []byte
	conn    *websocket.Conn
	decoder chan Decoder
	reason  string
}

func (s *Subscribe) Run(ctx context.Context) {
	if err := s.Init(ctx); err != nil {
		s.reason = err.Error()
		s.decoder <- &SubscribeMessage{
			Payload: gqlResponseErrBytes(GQLConnectionInit, err.Error()),
		}
		return
	}
	for {
		var msg SubscribeMessage
		if err := wsjson.Read(ctx, s.conn, &msg); err != nil {
			msg.Payload = gqlResponseErrBytes(GQLInvalidMessage, err.Error())
			s.decoder <- &msg
			return
		}
		switch msg.Type {
		case GQLConnectionAck:
			if err := s.subscribe(ctx); err != nil {
				msg.Payload = gqlResponseErrBytes(GQLSubscribe, err.Error())
				s.decoder <- &msg
				return
			}
		case GQLPing:
			if err := s.pong(ctx); err != nil {
				msg.Payload = gqlResponseErrBytes(GQLPong, err.Error())
				s.decoder <- &msg
				return
			}
		case GQLNext, GQLError:
			s.decoder <- &msg
		case GQLComplete:
			msg.Payload = gqlResponseErrBytes(GQLComplete, "server completed")
			s.decoder <- &msg
			return
		default:
			msg.Payload = gqlResponseErrBytes(msg.Type, "unknown message type")
			s.decoder <- &msg
			return
		}
	}
}

func (s *Subscribe) Init(ctx context.Context) error {
	fmt.Println(s.url.String())
	var err error
	s.conn, _, err = websocket.Dial(ctx, s.url.String(), &websocket.DialOptions{Subprotocols: []string{"graphql-transport-ws"}})
	if err != nil {
		return err
	}

	msg := &SubscribeMessage{
		Type: GQLConnectionInit,
	}
	var bParams []byte = nil
	if s.header != nil {
		params := map[string]interface{}{"headers": s.header}
		bParams, err = json.Marshal(params)
		if err != nil {
			return err
		}
		msg.Payload = bParams
	}

	return s.send(ctx, msg)
}

func (s *Subscribe) send(ctx context.Context, msg *SubscribeMessage) error {
	return wsjson.Write(ctx, s.conn, msg)
}

func (s *Subscribe) subscribe(ctx context.Context) error {
	subscribe := SubscribeMessage{
		ID:      uuid.New().String(),
		Type:    GQLSubscribe,
		Payload: s.payload,
	}
	return wsjson.Write(ctx, s.conn, subscribe)
}

func (s *Subscribe) pong(ctx context.Context) error {
	pong := SubscribeMessage{
		Type: GQLPong,
	}
	return wsjson.Write(ctx, s.conn, pong)
}

func (s *Subscribe) message(ctx context.Context, ch chan<- Event) {
	var event SubscribeEvent
	if err := wsjson.Read(ctx, s.conn, event.message); err != nil {
		event.error = &SubscribeError{
			Code:   websocket.StatusUnsupportedData,
			Reason: err.Error(),
		}
		ch <- &event
		return
	}
	switch event.message.Type {
	case GQLConnectionAck:
		if err := s.subscribe(ctx); err != nil {
			event.error = &SubscribeError{
				Code:   websocket.StatusInvalidFramePayloadData,
				Reason: err.Error(),
			}
			ch <- &event
		}
	case GQLPing:
		if err := s.pong(ctx); err != nil {
			event.error = &SubscribeError{
				Code:   websocket.StatusBadGateway,
				Reason: err.Error(),
			}
			ch <- &event
		}
	case GQLNext:
		ch <- &event
	case GQLError:
		event.error = &SubscribeError{
			Code:   websocket.StatusInternalError,
			Reason: "operation execution error(s) in response to the Subscribe message",
		}
		ch <- &event
	case GQLComplete:
		event.error = &SubscribeError{
			Code:   websocket.StatusNormalClosure,
			Reason: "operation execution has completed",
		}
		ch <- &event
	default:
		event.error = &SubscribeError{
			Code:   websocket.StatusNoStatusRcvd,
			Reason: "no status reserved",
		}
		ch <- &event
	}
}

func (s *Subscribe) Stop() {
	if s.conn != nil {
		s.conn.Close(websocket.StatusNormalClosure, s.reason)
	}
	close(s.decoder)
}

func (s *Subscribe) ResultChan() <-chan Decoder {
	return s.decoder
}
