package notify

import (
	"io/ioutil"
	"path"

	reference2 "github.com/floppyzedolfin/sync/reference/client/reference"
	"github.com/gofiber/fiber/v2"
)

const (
	// this path is mounted by the docker to the location specified at runtime
	localReplicaRoot = "/data/replica"
)

// Notify copies the content of the file locally
func Notify(_ *fiber.Ctx,  req *reference2.NotifyRequest) (*reference2.NotifyResponse, *fiber.Error) {
	err := ioutil.WriteFile(path.Join(localReplicaRoot, req.FullPath), []byte(req.FullContents), 0755)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError)
	}
	return &reference2.NotifyResponse{}, nil
}
