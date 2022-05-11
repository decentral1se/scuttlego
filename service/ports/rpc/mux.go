// Package rpc implements handlers for incoming Secure Scuttlebutt RPC requests.
package rpc

import (
	"github.com/planetary-social/go-ssb/service/domain/transport/rpc/mux"
)

// NewMuxHandlers is a convenience function used to create a list of all
// handlers implemented by this program.
func NewMuxHandlers(
	createHistoryStream *HandlerCreateHistoryStream,
	blobsGet *HandlerBlobsGet,
) []mux.Handler {
	return []mux.Handler{
		createHistoryStream,
		blobsGet,
	}
}
