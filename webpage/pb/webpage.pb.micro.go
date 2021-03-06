// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: webpage.proto

package pb

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

// Api Endpoints for WebpageService service

func NewWebpageServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for WebpageService service

type WebpageService interface {
	NewWebpage(ctx context.Context, in *WebpageRequest, opts ...client.CallOption) (*WebpageResponse, error)
}

type webpageService struct {
	c    client.Client
	name string
}

func NewWebpageService(name string, c client.Client) WebpageService {
	return &webpageService{
		c:    c,
		name: name,
	}
}

func (c *webpageService) NewWebpage(ctx context.Context, in *WebpageRequest, opts ...client.CallOption) (*WebpageResponse, error) {
	req := c.c.NewRequest(c.name, "WebpageService.NewWebpage", in)
	out := new(WebpageResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for WebpageService service

type WebpageServiceHandler interface {
	NewWebpage(context.Context, *WebpageRequest, *WebpageResponse) error
}

func RegisterWebpageServiceHandler(s server.Server, hdlr WebpageServiceHandler, opts ...server.HandlerOption) error {
	type webpageService interface {
		NewWebpage(ctx context.Context, in *WebpageRequest, out *WebpageResponse) error
	}
	type WebpageService struct {
		webpageService
	}
	h := &webpageServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&WebpageService{h}, opts...))
}

type webpageServiceHandler struct {
	WebpageServiceHandler
}

func (h *webpageServiceHandler) NewWebpage(ctx context.Context, in *WebpageRequest, out *WebpageResponse) error {
	return h.WebpageServiceHandler.NewWebpage(ctx, in, out)
}
