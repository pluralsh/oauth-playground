// based on ory's keto cli's grpc client:
// https://github.com/ory/keto/blob/6c0e1ba87f4d3a355cebd0ea77f28319be2dd606/cmd/client/grpc_client.go
package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/credentials/oauth"
)

type contextKeys string

const (
	ReadRemoteDefault  = "127.0.0.1:4466"
	WriteRemoteDefault = "127.0.0.1:4467"
	EnvReadRemote      = "KETO_READ_REMOTE"
	EnvWriteRemote     = "KETO_WRITE_REMOTE"
	EnvAuthToken       = "KETO_BEARER_TOKEN" // nosec G101 -- just the key, not the value
	EnvAuthority       = "KETO_AUTHORITY"

	ContextKeyTimeout contextKeys = "timeout"
)

type ConnectionDetails struct {
	readRemote, writeRemote, token, authority string
	skipHostVerification                      bool
	noTransportSecurity                       bool
}

func getRemote(envRemote, remoteDefault string) (remote string) {
	defer (func() {
		if strings.HasPrefix(remote, "http://") || strings.HasPrefix(remote, "https://") {
			_, _ = fmt.Fprintf(os.Stderr, "remote \"%s\" seems to be an http URL instead of a remote address\n", remote)
		}
	})()

	if remote, isSet := os.LookupEnv(envRemote); isSet {
		return remote
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "env var %s is not set, falling back to %s\n", envRemote, remoteDefault)
		return remoteDefault
	}
}

func (cd *ConnectionDetails) dialOptions() (opts []grpc.DialOption) {
	if cd.token != "" {
		opts = append(opts,
			grpc.WithPerRPCCredentials(
				oauth.NewOauthAccess(&oauth2.Token{AccessToken: cd.token})))
	}
	if cd.authority != "" {
		opts = append(opts, grpc.WithAuthority(cd.authority))
	}

	// TLS settings
	switch {
	case cd.noTransportSecurity:
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	case cd.skipHostVerification:
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			// nolint explicity set through scary flag
			InsecureSkipVerify: true,
		})))
	default:
		// Defaults to the default host root CA bundle
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(nil)))
	}
	return opts
}

func getAuthority() string {
	return os.Getenv(EnvAuthority)
}

func NewConnectionDetailsFromEnv() ConnectionDetails {
	return ConnectionDetails{
		readRemote:           getRemote(EnvReadRemote, ReadRemoteDefault),
		writeRemote:          getRemote(EnvWriteRemote, WriteRemoteDefault),
		token:                os.Getenv(EnvAuthToken),
		authority:            getAuthority(),
		skipHostVerification: true,
		noTransportSecurity:  true,
	}
}

func (cd *ConnectionDetails) ReadConn(ctx context.Context) (*grpc.ClientConn, error) {
	return Conn(ctx,
		cd.readRemote,
		cd,
	)
}

func (cd *ConnectionDetails) WriteConn(ctx context.Context) (*grpc.ClientConn, error) {
	return Conn(ctx,
		cd.writeRemote,
		cd,
	)
}

func Conn(ctx context.Context, remote string, cd *ConnectionDetails) (*grpc.ClientConn, error) {
	timeout := 3 * time.Second
	if d, ok := ctx.Value(ContextKeyTimeout).(time.Duration); ok {
		timeout = d
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return grpc.DialContext(
		ctx,
		remote,
		append([]grpc.DialOption{
			grpc.WithBlock(),
			grpc.WithDisableHealthCheck(),
		}, cd.dialOptions()...)...,
	)
}
