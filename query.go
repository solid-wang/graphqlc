package graphqlc

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type QueryOrMutate interface {
	Do(ctx context.Context) Decoder
}

type Query struct {
	url    *url.URL
	header map[string]string
	req    io.Reader
}

func (q *Query) Do(ctx context.Context) Decoder {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, q.url.String(), q.req)
	if err != nil {
		return &GraphResponse{Errors: []GraphError{{Message: err.Error()}}}
	}
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	for k, v := range q.header {
		request.Header.Add(k, v)
	}
	result, err := http.DefaultClient.Do(request)
	if err != nil {
		return &GraphResponse{Errors: []GraphError{{Message: err.Error()}}}
	}
	data, err := io.ReadAll(result.Body)
	if err != nil {
		return &GraphResponse{Errors: []GraphError{{Message: err.Error()}}}
	}
	var resp GraphResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return &GraphResponse{Errors: []GraphError{{Message: err.Error()}}}
	}
	return &resp
}
