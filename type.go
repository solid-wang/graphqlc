package graphqlc

import (
	"encoding/json"
	"fmt"
	"nhooyr.io/websocket"
)

// GraphRequest is a standard GraphQL POST request
// https://graphql.org/learn/serving-over-http/#post-request
type GraphRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
	// OperationName is only required if multiple operations are present in the query.
	OperationName string `json:"operationName,omitempty"`
}

func NewGraphRequest(query string, variables map[string]any) *GraphRequest {
	return &GraphRequest{
		Query:     query,
		Variables: variables,
	}
}

type Decoder interface {
	Decode(v interface{}) error
}

// GraphResponse is a standard GraphQL response
// https://graphql.org/learn/serving-over-http/#response
type GraphResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []GraphError    `json:"errors"`
}

type GraphError struct {
	Message string `json:"message"`
}

func (ge GraphError) Error() string {
	return ge.Message
}

func (gr *GraphResponse) Decode(v interface{}) error {
	if len(gr.Errors) > 0 {
		return gr.Errors[0]
	}
	return json.Unmarshal(gr.Data, v)
}

// SubscribeMessage represents a subscription operation message
// graphql-ws: https://github.com/enisdenjo/graphql-ws/blob/master/PROTOCOL.md
type SubscribeMessage struct {
	ID      string               `json:"id,omitempty"`
	Type    SubscribeMessageType `json:"type"`
	Payload json.RawMessage      `json:"payload,omitempty"`
}

func (sm *SubscribeMessage) Decode(v interface{}) error {
	var gr GraphResponse
	err := json.Unmarshal(sm.Payload, &gr)
	if err != nil {
		return err
	}
	return gr.Decode(v)
}

// SubscribeMessageType represents a subscription message enum type
type SubscribeMessageType string

const (
	// GQLConnectionInit Direction: Client -> Server
	// Indicates that the client wants to establish a connection within the existing socket.
	// This connection is not the actual WebSocket communication channel, but is rather a frame within it asking the server to allow future operation requests.
	GQLConnectionInit SubscribeMessageType = "connection_init"

	// GQLConnectionAck Direction: Server -> Client
	// Expected response to the ConnectionInit message from the client acknowledging a successful connection with the server.
	GQLConnectionAck SubscribeMessageType = "connection_ack"

	// GQLPing Direction: bidirectional
	// The Ping message can be sent at any time within the established socket.
	GQLPing SubscribeMessageType = "ping"

	// GQLPong Direction: bidirectional
	// The response to the Ping message. Must be sent as soon as the Ping message is received.
	GQLPong SubscribeMessageType = "pong"

	// GQLSubscribe Direction: Client -> Server
	// Requests an operation specified in the message payload. This message provides a unique ID field to connect published messages to the operation requested by this message.
	GQLSubscribe SubscribeMessageType = "subscribe"

	// GQLNext Direction: Server -> Client
	// Operation execution result(s) from the source stream created by the binding Subscribe message. After all results have been emitted, the Complete message will follow indicating stream completion.
	GQLNext SubscribeMessageType = "next"

	// GQLError Direction: Server -> Client
	// Operation execution error(s) in response to the Subscribe message.
	// This can occur before execution starts, usually due to validation errors, or during the execution of the request.
	GQLError SubscribeMessageType = "error"

	// GQLComplete Direction: bidirectional
	// indicates that the requested operation execution has completed. If the server dispatched the Error message relative to the original Subscribe message, no Complete message will be emitted.
	GQLComplete SubscribeMessageType = "complete"

	// GQLInvalidMessage Direction: bidirectional
	// Receiving a message of a type or format which is not specified in this document will result in an immediate socket closure with the event 4400: <error-message>. The <error-message> can be vaguely descriptive on why the received message is invalid.
	GQLInvalidMessage SubscribeMessageType = "invalid message"
)

type SubscribeError struct {
	Code   websocket.StatusCode
	Reason string
}

func (s *SubscribeError) Error() string {
	return fmt.Errorf("%s %s", s.Code, s.Reason).Error()
}

type Event interface {
	Message() *SubscribeMessage
	Error() *SubscribeError
}

type SubscribeEvent struct {
	message *SubscribeMessage
	error   *SubscribeError
}

func (s *SubscribeEvent) Message() *SubscribeMessage {
	return s.message
}

func (s *SubscribeEvent) Error() *SubscribeError {
	return s.error
}
