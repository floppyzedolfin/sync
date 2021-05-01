package client

import (
	"context"

	"github.com/floppyzedolfin/sync/reference/client/reference"
)

// API is the exposed contract of the service
type API interface {
	Patch(context.Context, *reference.PatchRequest) (*reference.PatchResponse, error)
	Delete(context.Context, *reference.DeleteRequest) (*reference.DeleteResponse, error)
}
