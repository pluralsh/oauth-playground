package clients

import (
	"fmt"
	"os"

	kratos "github.com/ory/kratos-client-go"
)

const (
	KratosPublicDefault = "http://127.0.0.1:4433"
	KratosAdminDefault  = "http://127.0.0.1:4434"
	KratosEnvPublic     = "KRATOS_PUBLIC_URL"
	KratosEnvAdmin      = "KRATOS_ADMIN_URL"
)

func NewKratosAdminClient() (*kratos.APIClient, error) {
	kratosAdminUrl := os.Getenv(KratosEnvAdmin)
	if kratosAdminUrl == "" {
		return nil, fmt.Errorf("No admin address configured for kratos")
	}
	kratosAdminConfiguration := kratos.NewConfiguration()
	kratosAdminConfiguration.Servers = []kratos.ServerConfiguration{
		{
			URL: kratosAdminUrl, // Kratos Public API
		},
	}
	kratosAdminClient := kratos.NewAPIClient(kratosAdminConfiguration)
	return kratosAdminClient, nil
}

func NewKratosPublicClient() (*kratos.APIClient, error) {
	kratosPublicUrl := os.Getenv(KratosEnvPublic)
	if kratosPublicUrl == "" {
		return nil, fmt.Errorf("No public address configured for kratos")
	}
	kratosPublicConfiguration := kratos.NewConfiguration()
	kratosPublicConfiguration.Servers = []kratos.ServerConfiguration{
		{
			URL: kratosPublicUrl, // Kratos Public API
		},
	}
	kratosPublicClient := kratos.NewAPIClient(kratosPublicConfiguration)
	return kratosPublicClient, nil
}
