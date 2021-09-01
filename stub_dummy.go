//go:build (!windows && ignore) || (!linux && ignore) || (!amd64 && ignore) || !cgo
// +build !windows,ignore !linux,ignore !amd64,ignore !cgo

package fabric

import "unsafe"

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

func fabricGetLastError() string {
	return ""
}

func callfabricFabricCreateRuntime(
	clzid unsafe.Pointer,
	outptr unsafe.Pointer,
) error {
	return errComNotImpl
}

func callfabricFabricGetActivationContext(
	clzid unsafe.Pointer,
	outptr unsafe.Pointer,
) error {
	return errComNotImpl
}

func createCallback(cb interface{}) uintptr {
	panic("not impl")
}
