package delete

import (
	"fmt"
	"os"
	"path"

	reference "github.com/floppyzedolfin/sync/reference/client/reference"
	"github.com/gofiber/fiber/v2"
)

const (
	// this path is mounted by the docker to the location specified at runtime
	localReplicaRoot = "/data/replica"
)

// Delete removes a local file
func Delete(_ *fiber.Ctx,  req *reference.DeleteRequest) (*reference.DeleteResponse, *fiber.Error) {
	fmt.Println("received request to delete " + req.FullPath)
	err := os.Remove(path.Join(localReplicaRoot, req.FullPath))
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError)
	}
	return &reference.DeleteResponse{}, nil
}
