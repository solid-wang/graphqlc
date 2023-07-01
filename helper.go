package graphqlc

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func encode(v interface{}) []byte {
	var reqBuf bytes.Buffer
	json.NewEncoder(&reqBuf).Encode(v)
	return reqBuf.Bytes()
}

func gqlResponseErrBytes(t SubscribeMessageType, reason string) []byte {
	return encode(GraphResponse{
		Errors: []GraphError{
			{
				Message: fmt.Sprintf("%s: %s", t, reason),
			},
		},
	})
}
