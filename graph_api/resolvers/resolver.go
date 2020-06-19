package resolvers

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	apiHolder "tutorial/graph_api/api_holder"
)

type Resolver struct {
	apiHolder apiHolder.Services
}

func NewResolver(ah apiHolder.Services) *Resolver {
	return &Resolver{apiHolder: ah}
}
