package handlers

import "github.com/pluralsh/oauth-playground/api-server/clients"

type Handler struct {
	C *clients.ClientWrapper
}
