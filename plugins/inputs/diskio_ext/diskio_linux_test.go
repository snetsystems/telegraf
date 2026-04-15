//go:build linux

package diskio_ext

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiskInfo(t *testing.T) {
	plugin := &DiskIO{
		infoCache: map[string]diskInfoCache{
			"null": {
				modifiedAt:   0,
				udevDataPath: "testdata/udev.txt",
				sysBlockPath: "testdata",
				values:       map[string]string{},
			},
		},
	}

	di, err := plugin.diskInfo("null")
	require.NoError(t, err)
	require.Equal(t, "myval1", di["MY_PARAM_1"])
	require.Equal(t, "myval2", di["MY_PARAM_2"])
	require.Equal(t, "/dev/foo/bar/devlink /dev/foo/bar/devlink1", di["DEVLINKS"])
}

// DiskIOStats.diskName isn't a linux specific function, but dependent
// functions are a no-op on non-Linux.
func TestDiskIOStats_diskName(t *testing.T) {
	tests := []struct {
		templates []string
		expected  string
	}{
		{[]string{"$MY_PARAM_1"}, "myval1"},
		{[]string{"${MY_PARAM_1}"}, "myval1"},
		{[]string{"x$MY_PARAM_1"}, "xmyval1"},
		{[]string{"x${MY_PARAM_1}x"}, "xmyval1x"},
		{[]string{"$MISSING", "$MY_PARAM_1"}, "myval1"},
		{[]string{"$MY_PARAM_1", "$MY_PARAM_2"}, "myval1"},
		{[]string{"$MISSING"}, "null"},
		{[]string{"$MY_PARAM_1/$MY_PARAM_2"}, "myval1/myval2"},
		{[]string{"$MY_PARAM_2/$MISSING"}, "null"},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("template %d", i), func(t *testing.T) {
			plugin := DiskIO{
				NameTemplates: tc.templates,
				infoCache: map[string]diskInfoCache{
					"null": {
						modifiedAt:   0,
						udevDataPath: "testdata/udev.txt",
						sysBlockPath: "testdata",
						values:       map[string]string{},
					},
				},
			}
			name, _ := plugin.diskName("null")
			require.Equal(t, tc.expected, name, "Templates: %#v", tc.templates)
		})
	}
}

// DiskIOStats.diskTags isn't a linux specific function, but dependent
// functions are a no-op on non-Linux.
func TestDiskIOStats_diskTags(t *testing.T) {
	plugin := &DiskIO{
		DeviceTags: []string{"MY_PARAM_2"},
		infoCache: map[string]diskInfoCache{
			"null": {
				modifiedAt:   0,
				udevDataPath: "testdata/udev.txt",
				sysBlockPath: "testdata",
				values:       map[string]string{},
			},
		},
	}
	dt := plugin.diskTags("null")
	require.Equal(t, map[string]string{"MY_PARAM_2": "myval2"}, dt)
}

func TestGetMountMap(t *testing.T) {
	plugin := &DiskIO{
		mountInfoPath: "testdata/mountinfo",
	}

	mountMap := plugin.getMountMap()

	// Verify basic parsing
	require.Equal(t, "/", mountMap["8:1"])
	require.Equal(t, "/home", mountMap["8:2"])
	require.Equal(t, "/var/log", mountMap["253:0"])

	// Verify first-wins: 8:1 appears twice (/ and /mnt/backup), should keep /
	require.Equal(t, "/", mountMap["8:1"])
}

func TestGetMountMap_FileNotFound(t *testing.T) {
	plugin := &DiskIO{
		mountInfoPath: "testdata/non-existent-mountinfo",
	}

	mountMap := plugin.getMountMap()

	// Should return empty map, not panic
	require.NotNil(t, mountMap)
	require.Empty(t, mountMap)
}

func TestGetDeviceMountPath_CacheHit(t *testing.T) {
	plugin := &DiskIO{
		infoCache: map[string]diskInfoCache{
			"sda1": {devID: "8:1"},
			"sda2": {devID: "8:2"},
		},
	}

	mountMap := map[string]string{
		"8:1": "/",
		"8:2": "/home",
	}

	require.Equal(t, "/", plugin.getDeviceMountPath("sda1", mountMap))
	require.Equal(t, "/home", plugin.getDeviceMountPath("sda2", mountMap))
	require.Equal(t, "", plugin.getDeviceMountPath("sda3", mountMap))
}

func TestGetDeviceMountPath_NilMap(t *testing.T) {
	plugin := &DiskIO{}
	require.Equal(t, "", plugin.getDeviceMountPath("sda1", nil))
}
