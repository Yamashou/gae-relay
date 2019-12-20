package main

import (
	"context"
	"github.com/99designs/gqlgen/handler"
	"net/http"

	"github.com/Yamashou/gae-relay/datastore"

	"go.mercari.io/datastore/aedatastore"

	"github.com/Yamashou/gae-relay/gqlgen"
	"github.com/Yamashou/gae-relay/server"

	"google.golang.org/appengine"
)

func main() {
	ctx := context.Background()
	datastoreClient, err := aedatastore.FromContext(ctx)
	if err != nil {
		panic(err)
	}
	schema := gqlgen.NewExecutableSchema(
		gqlgen.Config{Resolvers: server.NewResolver(datastore.NewClient(datastoreClient))},
	)
	http.Handle("/", handler.Playground("GAE Simple example", "/query"))
	http.Handle("/query", handler.GraphQL(
		schema,
	))
	appengine.Main()
}
