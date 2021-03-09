package fabric

import (
	"unsafe"

	ole "github.com/go-ole/go-ole"
	"golang.org/x/sys/windows"
)

var (
	fabricRuntimeDll              = windows.MustLoadDLL("FabricRuntime.dll")
	fabricFabricCreateRuntimeProc = fabricRuntimeDll.MustFindProc("FabricCreateRuntime")
	// TODO support this
	fabricFabricGetActivationContextProc = fabricRuntimeDll.MustFindProc("FabricGetActivationContext")
)

type FabricRuntime struct {
	hub *fabricRuntimeComHub
}

func NewFabricRuntime() (*FabricRuntime, error) {
	var com *comFabricRuntime

	clzid, _ := ole.IIDFromString("{cc53af8e-74cd-11df-ac3e-0024811e3892}")
	r, _, err := fabricCreateLocalClientProc.Call(uintptr(unsafe.Pointer(clzid)), uintptr(unsafe.Pointer(&com)))

	if r != 0 {
		return nil, err
	}

	// no init needed here, single com interface
	return &FabricRuntime{
		hub: &fabricRuntimeComHub{
			FabricRuntime: com,
		},
	}, nil
}
