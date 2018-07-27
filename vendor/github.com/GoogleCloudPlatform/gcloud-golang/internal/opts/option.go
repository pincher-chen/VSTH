// Package opts holds the DialOpts struct, configurable by
// cloud.ClientOptions to set up transports for cloud packages.
//
// This is a separate page to prevent cycles between the core
// cloud packages.
package opts

import (
	"net/http"

	"github.com/golang/oauth2"
	"github.com/grpc/grpc-go"
)

type DialOpt struct {
	Endpoint  string
	Scopes    []string
	UserAgent string

	TokenSource oauth2.TokenSource

	HTTPClient   *http.Client
	GRPCClient   *grpc.ClientConn
	GRPCDialOpts []grpc.DialOption
}
