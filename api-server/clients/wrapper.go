package clients

import (
	"github.com/go-logr/logr"
	hydra "github.com/ory/hydra-client-go/v2"
	kratos "github.com/ory/kratos-client-go"
	"github.com/pluralsh/oauth-playground/api-server/ent"
)

type ClientWrapper struct {
	KratosClient *kratos.APIClient
	KetoClient   *KetoGrpcClient
	HydraClient  *hydra.APIClient
	DbClient     *ent.Client
	Log          logr.Logger
}
