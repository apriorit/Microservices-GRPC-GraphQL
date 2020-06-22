package resolvers

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"tutorial/graph_api/services"
)

type Resolver struct {
	services services.Services
}

func NewResolver(s services.Services) *Resolver {
	return &Resolver{services: s}
}
