// +build windows

package fabric

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	fabricClientDll              = windows.NewLazyDLL("FabricClient.dll")
	fabricCreateLocalClientProc  = fabricClientDll.NewProc("FabricCreateLocalClient")
	fabricCreateLocalClient2Proc = fabricClientDll.NewProc("FabricCreateLocalClient2")
	fabricCreateLocalClient3Proc = fabricClientDll.NewProc("FabricCreateLocalClient3")
	fabricCreateLocalClient4Proc = fabricClientDll.NewProc("FabricCreateLocalClient4")
	fabricCreateClientProc       = fabricClientDll.NewProc("FabricCreateClient")
	fabricCreateClient2Proc      = fabricClientDll.NewProc("FabricCreateClient2")
	fabricCreateClient3Proc      = fabricClientDll.NewProc("FabricCreateClient3")
)

func callCreateClient3(
	len int,
	conns unsafe.Pointer,
	notificationHandler unsafe.Pointer,
	connectionHandler unsafe.Pointer,
	clzid unsafe.Pointer,
	outptr unsafe.Pointer,
) error {
	r, _, err := fabricCreateClient3Proc.Call(
		uintptr(len),
		uintptr(conns),
		uintptr(notificationHandler),
		uintptr(connectionHandler),
		uintptr(clzid),
		uintptr(outptr),
	)

	if r != 0 {
		return err
	}

	return nil
}

func callCreateLocalClient4(
	notificationHandler unsafe.Pointer,
	connectionHandler unsafe.Pointer,
	role FabricClientRole,
	clzid unsafe.Pointer,
	outptr unsafe.Pointer,
) error {
	r, _, err := fabricCreateLocalClient4Proc.Call(
		uintptr(notificationHandler),
		uintptr(connectionHandler),
		uintptr(role),
		uintptr(clzid),
		uintptr(outptr),
	)

	if r != 0 {
		return err
	}

	return nil
}
