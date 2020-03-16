package trash

import (
	// "fmt"
	// "math"
	// "strconv"
	// "strings"
	// "sync"
	// "time"

	// "github.com/juju/errors"

	"xsky-demon/log"
	"xsky-demon/models"
)

// VolumeWorker implements tash.Worker
type VolumeWorker struct {
	// mgmt *trash.trashMgmt
}

var _ Worker = (*VolumeWorker)(nil)

// Recycle moves volume resource to trash
func (worker *VolumeWorker) Recycle(res *models.TrashResource) error {
	return nil
}

// Restore restore volume resource
func (worker *VolumeWorker) Restore(res *models.TrashResource) error {
	return nil
}

// Delete delete volume resource from trash
func (worker *VolumeWorker) Delete(res *models.TrashResource) error {
	log.Infof("Delete resource %v from recycle bin", res)
	// create task
	// create acton log
	return nil
}

// UpdatePolicy modify the volume trash policy
func (worker *VolumeWorker) UpdatePolicy(policy *models.TrashPolicy) error {
	return nil
}

func init() {
	RegisterWorker(models.ResourceBlockVolume, new(VolumeWorker))
}
