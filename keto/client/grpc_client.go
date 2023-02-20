//based on ory's keto cli's grpc client:
//https://github.com/ory/keto/blob/6c0e1ba87f4d3a355cebd0ea77f28319be2dd606/cmd/client/grpc_client.go

package client

import (
	"context"
	. "context"

	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
	"google.golang.org/grpc"
	grpcHealthV1 "google.golang.org/grpc/health/grpc_health_v1"
)

type Client interface {
	//TODO
	//queryNamespaces()
	TransactTuples(ins []*rts.RelationTuple, del []*rts.RelationTuple)
	CreateTuple(r *rts.RelationTuple) error
	DeleteTuple(r *rts.RelationTuple) error
	DeleteAllTuples(q *rts.RelationQuery) error
	QueryTuple(q *rts.RelationQuery, opts ...PaginationOptionSetter) (*rts.ListRelationTuplesResponse, error)
	QueryAllTuples(q *rts.RelationQuery, pagesize int) ([]*rts.RelationTuple, error)
	Check(r *rts.RelationTuple) (error, bool)
	Expand(r *rts.SubjectSet, depth int) (error, *rts.SubjectTree)
	WaitUntilLive()
}

type GrpcClient struct {
	ConnDetails ConnectionDetails
	wc, rc      *grpc.ClientConn
	ctx         Context
}

func NewGrpcClient(ctx Context, cd ConnectionDetails) (*GrpcClient, error) {
	grpcClient := &GrpcClient{
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

func (g *GrpcClient) TransactTuples(ctx Context, ins []*rts.RelationTuple, del []*rts.RelationTuple) error {
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

func (g *GrpcClient) CreateTuple(ctx Context, r *rts.RelationTuple) error {
	return g.TransactTuples(ctx, []*rts.RelationTuple{r}, nil)
}

func (g *GrpcClient) CreateTuples(ctx Context, r []*rts.RelationTuple) error {
	return g.TransactTuples(ctx, r, nil)
}

func (g *GrpcClient) DeleteTuple(ctx Context, r *rts.RelationTuple) error {
	return g.TransactTuples(ctx, nil, []*rts.RelationTuple{r})
}

func (g *GrpcClient) DeleteTuples(ctx Context, r []*rts.RelationTuple) error {
	return g.TransactTuples(ctx, nil, r)
}

func (g *GrpcClient) DeleteAllTuples(ctx Context, q *rts.RelationQuery) error {
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

func WithToken(t string) PaginationOptionSetter {
	return func(opts *PaginationOptions) *PaginationOptions {
		opts.Token = t
		return opts
	}
}

func WithSize(size int) PaginationOptionSetter {
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

func (g *GrpcClient) QueryTuple(ctx Context, q *rts.RelationQuery, opts ...PaginationOptionSetter) (*rts.ListRelationTuplesResponse, error) {
	c := rts.NewReadServiceClient(g.rc)
	pagination := GetPaginationOptions(opts...)
	resp, err := c.ListRelationTuples(ctx, &rts.ListRelationTuplesRequest{
		RelationQuery: q,
		PageToken:     pagination.Token,
		PageSize:      int32(pagination.Size),
	})
	return resp, err
}

func (g *GrpcClient) QueryAllTuples(ctx Context, q *rts.RelationQuery, pagesize int) ([]*rts.RelationTuple, error) {
	tuples := make([]*rts.RelationTuple, 0)
	resp, err := g.QueryTuple(ctx, q, WithSize(pagesize))
	tuples = append(tuples, resp.RelationTuples...)
	for resp.NextPageToken != "" && err != nil {
		resp, err = g.QueryTuple(ctx, q, WithToken(resp.NextPageToken), WithSize(pagesize))
		tuples = append(tuples, resp.RelationTuples...)
	}
	return tuples, err
}

func (g *GrpcClient) Check(ctx Context, r *rts.RelationTuple) (bool, error) {
	c := rts.NewCheckServiceClient(g.rc)

	req := &rts.CheckRequest{
		Tuple: r,
	}
	resp, err := c.Check(ctx, req)

	return resp.Allowed, err
}

func (g *GrpcClient) Expand(ctx Context, ss *rts.Subject, depth int) (*rts.SubjectTree, error) {
	c := rts.NewExpandServiceClient(g.rc)

	resp, err := c.Expand(ctx, &rts.ExpandRequest{
		Subject:  ss,
		MaxDepth: int32(depth),
	})
	return resp.Tree, err
}

// TODO: not sure if this is the correct thing to do
func (g *GrpcClient) WaitUntilLive(ctx Context) error {
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
