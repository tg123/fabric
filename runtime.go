package fabric

import (
	"time"
	"unsafe"

	ole "github.com/go-ole/go-ole"
)

type FabricRuntime struct {
	coclz
	hub *fabricRuntimeComHub
}

var (
	iidIFabricRuntime = ole.NewGUID("{CC53AF8E-74CD-11DF-AC3E-0024811E3892}")
)

func NewFabricRuntime() (*FabricRuntime, error) {
	var com *comFabricRuntime

	err := callfabricFabricCreateRuntime(unsafe.Pointer(iidIFabricRuntime), unsafe.Pointer(&com))
	if err != nil {
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

	err := callfabricFabricGetActivationContext(
		unsafe.Pointer(ole.IID_IUnknown),
		unsafe.Pointer(&unknown),
	)

	if err != nil {
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
