package fabric

import (
	"time"
	"unsafe"

	ole "github.com/go-ole/go-ole"
	"golang.org/x/sys/windows"
)

var (
	fabricRuntimeDll              = windows.NewLazyDLL("FabricRuntime.dll")
	fabricFabricCreateRuntimeProc = fabricRuntimeDll.NewProc("FabricCreateRuntime")
	// TODO support this
	fabricFabricGetActivationContextProc = fabricRuntimeDll.NewProc("FabricGetActivationContext")
)

type FabricRuntime struct {
	hub            *fabricRuntimeComHub
	defaultTimeout time.Duration
}

func (v *FabricRuntime) GetTimeout() time.Duration {
	return v.defaultTimeout
}

func (v *FabricRuntime) SetDefaultTimeout(t time.Duration) {
	v.defaultTimeout = t
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
		defaultTimeout: 5 * time.Minute,
	}, nil
}
