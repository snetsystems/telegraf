//go:build !linux

package diskio_ext

type diskInfoCache struct{}

func (*DiskIO) diskInfo(_ string) (map[string]string, error) {
	return nil, nil
}

func resolveName(name string) string {
	return name
}

func getDeviceWWID(_ string) string {
	return ""
}

func (d *DiskIO) getMountMap() map[string]string {
	return nil
}

func (d *DiskIO) getDeviceMountPath(_ string, _ map[string]string) string {
	return ""
}
