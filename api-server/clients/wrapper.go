package clients

import (
	"github.com/go-logr/logr"
	hydra "github.com/ory/hydra-client-go/v2"
	kratos "github.com/ory/kratos-client-go"
)

type ClientWrapper struct {
	KratosAdminClient  *kratos.APIClient
	KratosPublicClient *kratos.APIClient
	KetoClient         *KetoGrpcClient
	HydraClient        *hydra.APIClient
	Log                logr.Logger
}
