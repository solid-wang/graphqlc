package graphqlc

import (
	"context"
	"testing"
	"time"
)

func TestSubscribe(t *testing.T) {
	var snowtooth = "http://snowtooth.moonhighway.com"
	client, err := NewClient(snowtooth)
	if err != nil {
		t.Errorf("new client err: %s", err)
	}
	req := NewGraphRequest(`
		subscription LiftStatusChange {
			liftStatusChange {
				id
				name
				status
				capacity
				night
				elevationGain
			}
		}
	`, nil)
	type LiftStatusChange struct {
		ID            string
		Name          string
		Status        string
		Capacity      int
		Night         bool
		ElevationGain int
	}
	type Response struct {
		LiftStatusChange LiftStatusChange
	}
	subscribe := client.Body(req).Subscription()
	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc()
	go subscribe.Run(ctx)
	defer subscribe.Stop()
	for {
		decoder := <-subscribe.ResultChan()
		var resp Response
		err := decoder.Decode(&resp)
		if err != nil {
			t.Error(err)
			return
		}
		t.Logf("get response: %v", resp)
	}
}
