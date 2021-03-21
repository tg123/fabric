package fabric

import (
	"time"
	"unsafe"

	ole "github.com/go-ole/go-ole"
	"golang.org/x/sys/windows"
)

var (
	fabricRuntimeDll                     = windows.NewLazyDLL("FabricRuntime.dll")
	fabricFabricCreateRuntimeProc        = fabricRuntimeDll.NewProc("FabricCreateRuntime")
	fabricFabricGetActivationContextProc = fabricRuntimeDll.NewProc("FabricGetActivationContext")
)

type FabricRuntime struct {
	coclz
	hub *fabricRuntimeComHub
}

func NewFabricRuntime() (*FabricRuntime, error) {
	var com *comFabricRuntime

	clzid, _ := ole.IIDFromString("{cc53af8e-74cd-11df-ac3e-0024811e3892}")
	r, _, err := fabricFabricCreateRuntimeProc.Call(uintptr(unsafe.Pointer(clzid)), uintptr(unsafe.Pointer(&com)))

	if r != 0 {
		return nil, err
	}

	hub := &fabricRuntimeComHub{
		FabricRuntime: com,
	}

	// no init needed here, single com interface
	return &FabricRuntime{
		coclz: coclz{
			hub:            hub,
			defaultTimeout: 5 * time.Minute,
		},
		hub: hub,
	}, nil
}

type FabricCodePackageActivationContext struct {
	coclz
	hub *fabricCodePackageActivationContextComHub
}

func NewFabricCodePackageActivationContext() (*FabricCodePackageActivationContext, error) {
	var unknown *ole.IUnknown

	r, _, err := fabricFabricGetActivationContextProc.Call(
		uintptr(unsafe.Pointer(ole.IID_IUnknown)),
		uintptr(unsafe.Pointer(&unknown)),
	)

	if r != 0 {
		return nil, err
	}

	hub := &fabricCodePackageActivationContextComHub{}
	hub.init(func(iid string, outptr unsafe.Pointer) error {
		return queryComObject(unknown, iid, outptr)
	})

	return &FabricCodePackageActivationContext{
		coclz: coclz{
			hub:            hub,
			defaultTimeout: 5 * time.Minute,
		},
		hub: hub,
	}, nil
}
