package shared

import (
	"github.com/hashicorp/go-plugin"
	"net/rpc"
)

type Greeter interface {
	Greet() string
}

type Foobar interface {
	Foobar() string
}

type GreeterRPC struct{ client *rpc.Client }

func (g *GreeterRPC) Greet() string {
	var resp string
	err := g.client.Call("Plugin.Greet", new(interface{}), &resp)
	if err != nil {
		panic(err)
	}

	return resp
}

type GreeterRPCServer struct {
	Impl Greeter
}

func (g *GreeterRPCServer) Greet(args interface{}, resp *string) error {
	*resp = g.Impl.Greet()
	return nil
}

type GreeterPlugin struct {
	Impl Greeter
}

func (p *GreeterPlugin) Server(broker *plugin.MuxBroker) (interface{}, error) {
	return &GreeterRPCServer{Impl: p.Impl}, nil
}

func (p *GreeterPlugin) Client(broker *plugin.MuxBroker, client *rpc.Client) (interface{}, error) {
	return &GreeterRPC{client: client}, nil
}

///

type FoobarRPC struct{ client *rpc.Client }

func (g *FoobarRPC) Foobar() string {
	var resp string
	err := g.client.Call("Plugin.Foobar", new(interface{}), &resp)
	if err != nil {
		panic(err)
	}

	return resp
}

type FoobarRPCServer struct {
	Impl Foobar
}

func (g *FoobarRPCServer) Foobar(args interface{}, resp *string) error {
	*resp = g.Impl.Foobar()
	return nil
}

type FoobarPlugin struct {
	Impl Foobar
}

func (p *FoobarPlugin) Server(broker *plugin.MuxBroker) (interface{}, error) {
	return &FoobarRPCServer{Impl: p.Impl}, nil
}

func (p *FoobarPlugin) Client(broker *plugin.MuxBroker, client *rpc.Client) (interface{}, error) {
	return &FoobarRPC{client: client}, nil
}

type BothRPCServer struct {
	Impl       Greeter
	FoobarImpl Foobar
}

func (g *BothRPCServer) Greet(args interface{}, resp *string) error {
	*resp = g.Impl.Greet()
	return nil
}

func (g *BothRPCServer) Foobar(args interface{}, resp *string) error {
	*resp = g.FoobarImpl.Foobar()
	return nil
}
