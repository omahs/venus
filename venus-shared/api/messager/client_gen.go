// Code generated by github.com/filecoin-project/venus/venus-devtool/api-gen. DO NOT EDIT.
package messager

import (
	"context"
	"net/http"

	"github.com/filecoin-project/go-jsonrpc"

	"github.com/filecoin-project/venus/venus-shared/api"
)

// NewIMessagerRPC creates a new httpparse jsonrpc remotecli.
func NewIMessagerRPC(ctx context.Context, addr string, requestHeader http.Header, opts ...jsonrpc.Option) (IMessager, jsonrpc.ClientCloser, error) {
	if requestHeader == nil {
		requestHeader = http.Header{}
	}
	requestHeader.Set(api.VenusAPINamespaceHeader, "messager.IMessager")

	var res IMessagerStruct
	closer, err := jsonrpc.NewMergeClient(ctx, addr, "Message", api.GetInternalStructs(&res), requestHeader, opts...)

	return &res, closer, err
}
