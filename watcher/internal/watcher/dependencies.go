package watcher

import (
	"context"

	"github.com/floppyzedolfin/sync/reference/client/reference"
)

type Reference interface {
	Patch(context.Context, reference.PatchRequest) (*reference.PatchResponse, error)
	Delete(context.Context, reference.DeleteRequest) (*reference.DeleteResponse, error)
}