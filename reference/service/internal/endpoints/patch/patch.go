package patch

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path"

	reference "github.com/floppyzedolfin/sync/reference/client/reference"
	"github.com/gofiber/fiber/v2"
)

const (
	// this path is mounted by the docker to the location specified at runtime
	localReplicaRoot = "/data/replica"
)

// Patch copies the content of the file locally
func Patch(_ *fiber.Ctx,  req *reference.PatchRequest) (*reference.PatchResponse, *fiber.Error) {
	fmt.Println("received request to patch " + req.FullPath)
	rights := getOctalRights(req.Rights)
	err := ioutil.WriteFile(path.Join(localReplicaRoot, req.FullPath), []byte(req.FullContents), fs.FileMode(rights))
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError)
	}
	return &reference.PatchResponse{}, nil
}

// getOctalRights converts 755 to 0755
func getOctalRights(decimalRights uint32) uint32 {
	return ((decimalRights / 100) & 0x7) << 2 + ((decimalRights / 10) & 0x7) << 1 + (decimalRights % 10) & 0x7
}