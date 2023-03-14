package handlers

import (
	"github.com/go-logr/logr"
	"github.com/pluralsh/oauth-playground/api-server/clients"
)

type Handler struct {
	C   *clients.ClientWrapper
	Log logr.Logger
}
