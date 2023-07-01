package graphqlc

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleQuery(t *testing.T) {
	var snowtooth = "http://snowtooth.moonhighway.com"
	client, err := NewClient(snowtooth)
	if err != nil {
		t.Errorf("new client err: %s", err)
	}
	req := NewGraphRequest(`
		query AllLifts {
			allLifts {
				id
				name
				status
				capacity
				night
				elevationGain
			}
		}
	`, nil)
	type Lifts struct {
		ID            string
		Name          string
		Status        string
		Capacity      int
		Night         bool
		ElevationGain int
	}
	type Response struct {
		AllLifts []Lifts
	}
	var resp Response
	err = client.Body(req).QueryOrMutate().Do(context.Background()).Decode(&resp)
	if err != nil {
		t.Errorf("test: %s", err)
	}
	fmt.Println(resp)
	t.Logf("get graphql response, lenth %d", len(resp.AllLifts))
}

func TestQueryWithVars(t *testing.T) {
	var snowtooth = "http://snowtooth.moonhighway.com"
	client, err := NewClient(snowtooth)
	if err != nil {
		t.Errorf("new client err: %s", err)
	}
	req := NewGraphRequest(`
		query AllLifts($status: LiftStatus) {
			allLifts(status: $status) {
				id
				name
				status
				capacity
				night
				elevationGain
			}
		}
	`, map[string]any{"status": "OPEN"})
	type Lifts struct {
		ID            string
		Name          string
		Status        string
		Capacity      int
		Night         bool
		ElevationGain int
	}
	type Response struct {
		AllLifts []Lifts
	}
	var resp Response
	err = client.Body(req).QueryOrMutate().Do(context.Background()).Decode(&resp)
	if err != nil {
		t.Errorf("test err: %s", err)
	}
	fmt.Println(resp)
	t.Logf("get graphql response, lenth %d", len(resp.AllLifts))
}

func TestQueryWithVarse(t *testing.T) {
	var snowtooth = "http://snowtooth.moonhighway.com"
	client, err := NewClient(snowtooth)
	if err != nil {
		t.Errorf("new client err: %s", err)
	}
	req := NewGraphRequest(`
		{
			allLifts(status: OPEN) {
				id
				name
				status
				capacity
				night
				elevationGain
			}
		}
	`, nil)
	type Lifts struct {
		ID            string
		Name          string
		Status        string
		Capacity      int
		Night         bool
		ElevationGain int
	}
	type Response struct {
		AllLifts []Lifts
	}
	var resp Response
	err = client.Body(req).QueryOrMutate().Do(context.Background()).Decode(&resp)
	if err != nil {
		t.Errorf("test err: %s", err)
	}
	fmt.Println(resp)
	t.Logf("get graphql response, lenth %d", len(resp.AllLifts))
}

func TestMultipleFields(t *testing.T) {
	var snowtooth = "http://snowtooth.moonhighway.com"
	client, err := NewClient(snowtooth)
	if err != nil {
		t.Errorf("new client err: %s", err)
	}
	req := NewGraphRequest(`
		query AllLifts($status: LiftStatus) {
			allLifts(status: $status) {
				name
				status
				capacity
				night
				elevationGain
				id
			}
			allTrails {
				id
				name
				status
				difficulty
				groomed
				trees
				night
			}
		}
	`, map[string]any{"status": "OPEN"})
	type Lifts struct {
		ID            string
		Name          string
		Status        string
		Capacity      int
		Night         bool
		ElevationGain int
	}
	type Trail struct {
		ID         string
		Name       string
		Status     string
		Difficulty string
		Groomed    bool
		Trees      bool
		Night      bool
	}
	type Response struct {
		AllLifts  []Lifts
		AllTrails []Trail
	}
	var resp Response
	err = client.Body(req).QueryOrMutate().Do(context.Background()).Decode(&resp)
	if err != nil {
		t.Errorf("test err: %s", err)
	}
	fmt.Println(resp)
	t.Logf("get lifts graphql response, lenth %d", len(resp.AllLifts))
	t.Logf("get trails graphql response, lenth %d", len(resp.AllTrails))
}

func TestMultipleOperations(t *testing.T) {
	var snowtooth = "http://snowtooth.moonhighway.com"
	client, err := NewClient(snowtooth)
	if err != nil {
		t.Errorf("new client err: %s", err)
	}
	req := NewGraphRequest(`
		query Lift($id: ID!) {
			Lift(id: $id) {
				id
				name
				status
				capacity
				night
				elevationGain
			}
		}
		mutation SetLiftStatus($id: ID!, $status: LiftStatus!) {
			setLiftStatus(id: $id, status: $status) {
				id
				name
				status
				capacity
				night
				elevationGain
			}
		}
	`, map[string]any{"id": "astra-express", "status": "OPEN"})
	type Lifts struct {
		ID            string
		Name          string
		Status        string
		Capacity      int
		Night         bool
		ElevationGain int
	}
	type SetResponse struct {
		SetLiftStatus Lifts
	}
	var setresp SetResponse
	err = client.Body(req).SetOperationName("SetLiftStatus").QueryOrMutate().Do(context.Background()).Decode(&setresp)
	if err != nil {
		t.Errorf("test: %s", err)
	}
	fmt.Println(setresp.SetLiftStatus)

	type LiftResponse struct {
		Lift Lifts
	}
	var liftresp LiftResponse
	err = client.Body(req).SetOperationName("Lift").QueryOrMutate().Do(context.Background()).Decode(&liftresp)
	if err != nil {
		t.Errorf("test: %s", err)
	}
	fmt.Println(liftresp.Lift)

	a := assert.New(t)
	a.Equal(setresp.SetLiftStatus, liftresp.Lift)
}
