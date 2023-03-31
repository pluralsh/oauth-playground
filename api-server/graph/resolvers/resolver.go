//go:generate go run generate.go
package resolvers

import "github.com/pluralsh/oauth-playground/api-server/clients"

type Resolver struct {
	C *clients.ClientWrapper
}
