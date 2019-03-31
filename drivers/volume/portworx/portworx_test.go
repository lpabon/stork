// This test requires an OpenStorage Server to be running. Here is how
// run it:
// docker run --rm --name sdk -d -p 9100:9100 -p 9110:9110 -p 9001:9001 openstorage/mock-sdk-server
//
// build unittest,mocksdkserver
package portworx

import (
	"testing"

	"github.com/stretchr/testify/assert"

	crdv1 "github.com/kubernetes-incubator/external-storage/snapshot/pkg/apis/crd/v1"
	v1 "k8s.io/api/core/v1"

	"github.com/libopenstorage/openstorage/api"
	volumeclient "github.com/libopenstorage/openstorage/api/client/volume"
	"github.com/libopenstorage/openstorage/volume"
)

func init() {
	driverName = "fake"
}

const (
	ostUrl = "http://localhost:9001"
)

func createVolume(t *testing.T, name string, v volume.VolumeDriver) string {
	size := uint64(1234)
	req := &api.VolumeCreateRequest{
		Locator: &api.VolumeLocator{
			Name: name,
			VolumeLabels: map[string]string{
				"hello": "world",
			},
		},
		Source: &api.Source{},
		Spec: &api.VolumeSpec{
			HaLevel: 3,
			Size:    size,
			Format:  api.FSType_FS_TYPE_EXT4,
			Shared:  true,
		},
	}
	id, err := v.Create(req.GetLocator(), req.GetSource(), req.GetSpec())
	assert.Nil(t, err)
	assert.NotEmpty(t, id)
	return id
}

func TestServerRunning(t *testing.T) {
	for i := 0; i < 5; i++ {
		// Create a volume using the REST interface
		cl, err := volumeclient.NewAuthDriverClient(ostUrl, "fake", "v1", "", "", "fake")
		assert.NoError(t, err)
		driverclient := volumeclient.VolumeDriver(cl)

		name := "myvol"
		id := createVolume(t, name, driverclient)

		vols, err := driverclient.Inspect([]string{id})
		assert.NoError(t, err)
		assert.NotNil(t, vols)
		assert.Len(t, vols, 1)
		assert.Equal(t, "myvol", vols[0].GetLocator().GetName())

		err = driverclient.Delete(id)
		assert.NoError(t, err)

		vols, err = driverclient.Inspect([]string{id})
		assert.NoError(t, err)
		assert.NotNil(t, vols)
		assert.Len(t, vols, 0)
	}
}

func TestInspectVolume(t *testing.T) {
	cl, err := volumeclient.NewAuthDriverClient(ostUrl, "fake", "v1", "", "", "fake")
	assert.NoError(t, err)
	driverclient := volumeclient.VolumeDriver(cl)

	p := &portworx{
		endpoint: ostUrl,
	}

	id := createVolume(t, "myvol", driverclient)
	defer driverclient.Delete(id)

	info, err := p.InspectVolume(id)
	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, id, info.VolumeID)
	assert.Equal(t, "myvol", info.VolumeName)
	assert.Len(t, info.Labels, 1)
	assert.Contains(t, info.Labels, "hello")
}

/*
func TestSnapshotCreate(t *testing.T) {
	cl, err := volumeclient.NewAuthDriverClient(ostUrl, "fake", "v1", "", "", "fake")
	assert.NoError(t, err)
	driverclient := volumeclient.VolumeDriver(cl)

	p := &portworx{
		endpoint: ostUrl,
	}

	id := createVolume(t, "myvol", driverclient)
	defer driverclient.Delete(id)

	snapcrd := &crdv1.VolumeSnapshot{}
	pv := &v1.PersistentVolume{
		Spec: v1.PersistentVolumeSpec{
			PersistemVolumeSource: v1.PersistentVolumeSource{
				PortworxVolume: &v1.PortworxVolumeSource{
					VolumeID: id,
				},
			},
		},
	}

}
*/
