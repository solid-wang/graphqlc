graphqlc
=======

Package `graphqlc` provides a GraphQL client implementation.

- [graphqlc](#graphqlc)
    - [Installation](#installation)
    - [Usage](#usage)
        - [Authentication](#authentication)
        - [Query](#query)
          - [Simple Query](#simple-query)
          - [Arguments and Variables](#arguments-and-variables)
        - [Mutations](#mutations)
        - [Subscription](#subscription)

## Installation

```bash
go get -u github.com/hasura/go-graphql-client
```

## Usage

Construct a GraphQL client, specifying the GraphQL server URL. Then, you can use it to make GraphQL queries、 mutations and subscription.

```Go
client, _ := graphqlc.NewClient("https://example.com/graphql")
// Use client...
```

### Query

#### Simple Query

For example, to make the following GraphQL query:

```GraphQL
query {
	me {
		name
	}
}
```

You can define this graphql request:

```Go
req := NewGraphRequest(`
    query Me {
        me {
            name
        }
    }
`, nil)
```

You also need to define the struct：
```Go
type Me struct {
    Name          string
}
type Response struct {
    Me Me
}
```

Then do request, you need passing a pointer to resp:

```Go
var resp Response
err := client.Body(req).QueryOrMutate().Do(context.Background()).Decode(&resp)
if err != nil {
	// Handle error.
}
fmt.Println(query.Me.Name)

// Output: Luke Skywalker
```

#### Arguments and Variables

Often, you'll want to specify arguments on some fields.

For example, to make the following GraphQL query:

```GraphQL
{
	human(id: "1000") {
		name
		height(unit: METER)
	}
}
```

You can define this graphql request:

```Go
req := NewGraphRequest(`
    {
        human(id: "1000") {
            name
            height(unit: METER)
        }
    }
`, nil)
```

Usually, you might want to use:
```Go
req := NewGraphRequest(`
    query Human($id: ID!, $unit: LengthUnit) {
        human(id: $id) {
            name
            height(unit: $unit)
        }
    }
`, map[string]any{"id": "1000", "unit": "METER"})
```

You also need to define the struct：
```Go
type Human struct {
    Name          string
	height          string
}
type Response struct {
    Human Human
}
```

Then do request:

```Go
var resp Response
err := client.Body(req).QueryOrMutate().Do(context.Background()).Decode(&resp)
if err != nil {
	// Handle error.
}
fmt.Println(q.Human.Name)
fmt.Println(q.Human.Height)

// Output:
// Luke Skywalker
// 1.72
```

### Mutations

Mutations often require information that you can only find out by performing a query first. Let's suppose you've already done that.

For example, to make the following GraphQL mutation:

```GraphQL
mutation($ep: Episode!, $review: ReviewInput!) {
	createReview(episode: $ep, review: $review) {
		stars
		commentary
	}
}
variables {
	"ep": "JEDI",
	"review": {
		"stars": 5,
		"commentary": "This is a great movie!"
	}
}
```

You can define:
```Go
req := NewGraphRequest(`
    mutation($ep: Episode!, $review: ReviewInput!) {
        createReview(episode: $ep, review: $review) {
            stars
            commentary
        }
    }
`, map[string]any{"ep": "JEDI", "review": map[string]any{"start": 5, "commentary": "This is a great movie!"}})
```

You also need to define the struct：
```Go
type CreateReview struct {
    stars          int
    commentary     string
}
type Response struct {
    CreateReview CreateReview
}
```

Then do request:

```Go
var resp Response
err = client.Body(req).QueryOrMutate().Do(context.Background()).Decode(&resp)
if err != nil {
	// Handle error.
}
fmt.Printf("Created a %v star review: %v\n", m.CreateReview.Stars, m.CreateReview.Commentary)

// Output:
// Created a 5 star review: This is a great movie!
```

### Subscription

For example, to make the following GraphQL query:

```GraphQL
subscription {
	me {
		name
	}
}
```

You can define this graphql request:

```Go
req := NewGraphRequest(`
    subscription {
        me {
            name
        }
    }
`, nil)
```

You also need to define the struct：
```Go
type Me struct {
    Name          string
}
type Response struct {
    Me Me
}
```

Then run subscribe, passing a pointer to it:

```Go
subscribe := client.Body(req).Subscription()
go subscribe.Run(context.Background())
defer subscribe.Stop()
for {
    decoder := <-subscribe.ResultChan()
    var resp Response
    err := decoder.Decode(&resp)
    if err != nil {
        // Handle error.
    return
    }
}
```
