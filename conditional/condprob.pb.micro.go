// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: condprob.proto

package conditional

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for ConditionalProbability service

func NewConditionalProbabilityEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for ConditionalProbability service

type ConditionalProbabilityService interface {
	CalcConditionalProbability(ctx context.Context, in *ConditionalProbabilityRequest, opts ...client.CallOption) (*ConditionalProbabilityResponse, error)
}

type conditionalProbabilityService struct {
	c    client.Client
	name string
}

func NewConditionalProbabilityService(name string, c client.Client) ConditionalProbabilityService {
	return &conditionalProbabilityService{
		c:    c,
		name: name,
	}
}

func (c *conditionalProbabilityService) CalcConditionalProbability(ctx context.Context, in *ConditionalProbabilityRequest, opts ...client.CallOption) (*ConditionalProbabilityResponse, error) {
	req := c.c.NewRequest(c.name, "ConditionalProbability.CalcConditionalProbability", in)
	out := new(ConditionalProbabilityResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ConditionalProbability service

type ConditionalProbabilityHandler interface {
	CalcConditionalProbability(context.Context, *ConditionalProbabilityRequest, *ConditionalProbabilityResponse) error
}

func RegisterConditionalProbabilityHandler(s server.Server, hdlr ConditionalProbabilityHandler, opts ...server.HandlerOption) error {
	type conditionalProbability interface {
		CalcConditionalProbability(ctx context.Context, in *ConditionalProbabilityRequest, out *ConditionalProbabilityResponse) error
	}
	type ConditionalProbability struct {
		conditionalProbability
	}
	h := &conditionalProbabilityHandler{hdlr}
	return s.Handle(s.NewHandler(&ConditionalProbability{h}, opts...))
}

type conditionalProbabilityHandler struct {
	ConditionalProbabilityHandler
}

func (h *conditionalProbabilityHandler) CalcConditionalProbability(ctx context.Context, in *ConditionalProbabilityRequest, out *ConditionalProbabilityResponse) error {
	return h.ConditionalProbabilityHandler.CalcConditionalProbability(ctx, in, out)
}