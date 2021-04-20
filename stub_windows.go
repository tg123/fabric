// +build windows, amd64

package fabric

import (
	"syscall"
	"unsafe"
)

var (
	fabricClientDll              = syscall.NewLazyDLL("FabricClient.dll")
	fabricCreateLocalClientProc  = fabricClientDll.NewProc("FabricCreateLocalClient")
	fabricCreateLocalClient2Proc = fabricClientDll.NewProc("FabricCreateLocalClient2")
	fabricCreateLocalClient3Proc = fabricClientDll.NewProc("FabricCreateLocalClient3")
	fabricCreateLocalClient4Proc = fabricClientDll.NewProc("FabricCreateLocalClient4")
	fabricCreateClientProc       = fabricClientDll.NewProc("FabricCreateClient")
	fabricCreateClient2Proc      = fabricClientDll.NewProc("FabricCreateClient2")
	fabricCreateClient3Proc      = fabricClientDll.NewProc("FabricCreateClient3")
)

func boolToUintptr(x bool) uintptr {
	if x {
		return 1
	}
	return 0
}

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

	return errno(r, err)
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

	return errno(r, err)
}

var (
	fabricCommonDll               = syscall.NewLazyDLL("FabricCommon.dll")
	fabricGetLastErrorMessageProc = fabricCommonDll.NewProc("FabricGetLastErrorMessage")
)

func fabricGetLastError() string {
	var result *comFabricStringResult
	hr, _, _ := fabricGetLastErrorMessageProc.Call(uintptr(unsafe.Pointer(&result)))

	if hr != 0 {
		return ""
	}

	msg, _ := result.GetString()
	return msg
}

var (
	fabricRuntimeDll                     = syscall.NewLazyDLL("FabricRuntime.dll")
	fabricFabricCreateRuntimeProc        = fabricRuntimeDll.NewProc("FabricCreateRuntime")
	fabricFabricGetActivationContextProc = fabricRuntimeDll.NewProc("FabricGetActivationContext")
)

func callfabricFabricCreateRuntime(
	clzid unsafe.Pointer,
	outptr unsafe.Pointer,
) error {
	r, _, err := fabricFabricCreateRuntimeProc.Call(uintptr(clzid), uintptr(outptr))
	return errno(r, err)
}

func callfabricFabricGetActivationContext(
	clzid unsafe.Pointer,
	outptr unsafe.Pointer,
) error {
	r, _, err := fabricFabricGetActivationContextProc.Call(uintptr(clzid), uintptr(outptr))
	return errno(r, err)
}

func createCallback(cb interface{}) uintptr {
	return syscall.NewCallback(cb)
}
