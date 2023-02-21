package clients

import (
	"context"
	. "context"
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

	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
	grpcHealthV1 "google.golang.org/grpc/health/grpc_health_v1"
)

type ketoContextKeys string

const (
	KetoReadRemoteDefault  = "127.0.0.1:4466"
	KetoWriteRemoteDefault = "127.0.0.1:4467"
	KetoEnvReadRemote      = "KETO_READ_REMOTE"
	KetoEnvWriteRemote     = "KETO_WRITE_REMOTE"
	KetoEnvAuthToken       = "KETO_BEARER_TOKEN" // nosec G101 -- just the key, not the value
	KetoEnvAuthority       = "KETO_AUTHORITY"

	ContextKeyTimeout ketoContextKeys = "timeout"
)

type KetoConnectionDetails struct {
	readRemote, writeRemote, token, authority string
	skipHostVerification                      bool
	noTransportSecurity                       bool
}

func getKetoRemote(envRemote, remoteDefault string) (remote string) {
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

func (cd *KetoConnectionDetails) dialOptions() (opts []grpc.DialOption) {
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

func getKetoAuthority() string {
	return os.Getenv(KetoEnvAuthority)
}

func NewKetoConnectionDetailsFromEnv() KetoConnectionDetails {
	return KetoConnectionDetails{
		readRemote:           getKetoRemote(KetoEnvReadRemote, KetoReadRemoteDefault),
		writeRemote:          getKetoRemote(KetoEnvWriteRemote, KetoWriteRemoteDefault),
		token:                os.Getenv(KetoEnvAuthToken),
		authority:            getKetoAuthority(),
		skipHostVerification: true,
		noTransportSecurity:  true,
	}
}

func (cd *KetoConnectionDetails) ReadConn(ctx context.Context) (*grpc.ClientConn, error) {
	return KetoConn(ctx,
		cd.readRemote,
		cd,
	)
}

func (cd *KetoConnectionDetails) WriteConn(ctx context.Context) (*grpc.ClientConn, error) {
	return KetoConn(ctx,
		cd.writeRemote,
		cd,
	)
}

func KetoConn(ctx context.Context, remote string, cd *KetoConnectionDetails) (*grpc.ClientConn, error) {
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

type KetoGrpcClient struct {
	ConnDetails KetoConnectionDetails
	wc, rc      *grpc.ClientConn
	ctx         Context
}

func NewKetoGrpcClient(ctx Context, cd KetoConnectionDetails) (*KetoGrpcClient, error) {
	grpcClient := &KetoGrpcClient{
		ConnDetails: cd,
		ctx:         ctx,
	}
	if wc, err := cd.WriteConn(ctx); err != nil {
		return nil, err
	} else {
		grpcClient.wc = wc
	}
	if rc, err := cd.ReadConn(ctx); err != nil {
		return nil, err
	} else {
		grpcClient.rc = rc
	}
	return grpcClient, nil
}

func (g *KetoGrpcClient) TransactTuples(ctx Context, ins []*rts.RelationTuple, del []*rts.RelationTuple) error {
	c := rts.NewWriteServiceClient(g.wc)

	deltas := append(
		rts.RelationTupleToDeltas(ins, rts.RelationTupleDelta_ACTION_INSERT),
		rts.RelationTupleToDeltas(del, rts.RelationTupleDelta_ACTION_DELETE)...,
	)

	_, err := c.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: deltas,
	})
	return err
}

func (g *KetoGrpcClient) CreateTuple(ctx Context, r *rts.RelationTuple) error {
	return g.TransactTuples(ctx, []*rts.RelationTuple{r}, nil)
}

func (g *KetoGrpcClient) CreateTuples(ctx Context, r []*rts.RelationTuple) error {
	return g.TransactTuples(ctx, r, nil)
}

func (g *KetoGrpcClient) DeleteTuple(ctx Context, r *rts.RelationTuple) error {
	return g.TransactTuples(ctx, nil, []*rts.RelationTuple{r})
}

func (g *KetoGrpcClient) DeleteTuples(ctx Context, r []*rts.RelationTuple) error {
	return g.TransactTuples(ctx, nil, r)
}

func (g *KetoGrpcClient) DeleteAllTuples(ctx Context, q *rts.RelationQuery) error {
	c := rts.NewWriteServiceClient(g.wc)
	_, err := c.DeleteRelationTuples(ctx, &rts.DeleteRelationTuplesRequest{
		RelationQuery: q,
	})
	return err
}

type (
	PaginationOptions struct {
		Token string `json:"page_token"`
		Size  int    `json:"page_size"`
	}
	PaginationOptionSetter func(*PaginationOptions) *PaginationOptions
)

func KetoWithToken(t string) PaginationOptionSetter {
	return func(opts *PaginationOptions) *PaginationOptions {
		opts.Token = t
		return opts
	}
}

func KetoWithSize(size int) PaginationOptionSetter {
	return func(opts *PaginationOptions) *PaginationOptions {
		opts.Size = size
		return opts
	}
}

func GetPaginationOptions(modifiers ...PaginationOptionSetter) *PaginationOptions {
	opts := &PaginationOptions{}
	for _, f := range modifiers {
		opts = f(opts)
	}
	return opts
}

func (g *KetoGrpcClient) QueryTuple(ctx Context, q *rts.RelationQuery, opts ...PaginationOptionSetter) (*rts.ListRelationTuplesResponse, error) {
	c := rts.NewReadServiceClient(g.rc)
	pagination := GetPaginationOptions(opts...)
	resp, err := c.ListRelationTuples(ctx, &rts.ListRelationTuplesRequest{
		RelationQuery: q,
		PageToken:     pagination.Token,
		PageSize:      int32(pagination.Size),
	})
	return resp, err
}

func (g *KetoGrpcClient) QueryAllTuples(ctx Context, q *rts.RelationQuery, pagesize int) ([]*rts.RelationTuple, error) {
	tuples := make([]*rts.RelationTuple, 0)
	resp, err := g.QueryTuple(ctx, q, KetoWithSize(pagesize))
	tuples = append(tuples, resp.RelationTuples...)
	for resp.NextPageToken != "" && err != nil {
		resp, err = g.QueryTuple(ctx, q, KetoWithToken(resp.NextPageToken), KetoWithSize(pagesize))
		tuples = append(tuples, resp.RelationTuples...)
	}
	return tuples, err
}

func (g *KetoGrpcClient) Check(ctx Context, r *rts.RelationTuple) (bool, error) {
	c := rts.NewCheckServiceClient(g.rc)

	req := &rts.CheckRequest{
		Tuple: r,
	}
	resp, err := c.Check(ctx, req)

	return resp.Allowed, err
}

func (g *KetoGrpcClient) Expand(ctx Context, ss *rts.Subject, depth int) (*rts.SubjectTree, error) {
	c := rts.NewExpandServiceClient(g.rc)

	resp, err := c.Expand(ctx, &rts.ExpandRequest{
		Subject:  ss,
		MaxDepth: int32(depth),
	})
	return resp.Tree, err
}

// TODO: not sure if this is the correct thing to do
func (g *KetoGrpcClient) WaitUntilLive(ctx Context) error {
	c := grpcHealthV1.NewHealthClient(g.rc)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cl, err := c.Watch(ctx, &grpcHealthV1.HealthCheckRequest{})
	if err != nil {
		return err
	}

	for {
		select {
		case <-g.ctx.Done():
			return nil
		default:
		}
		resp, err := cl.Recv()

		if resp.Status == grpcHealthV1.HealthCheckResponse_SERVING {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

// TODO: add namespace client
// func (g *KetoGrpcClient) Namespace() error {}
