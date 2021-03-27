// +build !windows

package fabric

func callCreateClient3(
	len int,
	conns unsafe.Pointer,
	notificationHandler unsafe.Pointer,
	connectionHandler unsafe.Pointer,
	clzid unsafe.Pointer,
	outptr unsafe.Pointer,
) error {
	return errComNotImpl
}

func callCreateLocalClient4(
	notificationHandler unsafe.Pointer,
	connectionHandler unsafe.Pointer,
	role FabricClientRole,
	clzid unsafe.Pointer,
	outptr unsafe.Pointer,
) error {
	return errComNotImpl
}
