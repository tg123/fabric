package fabric

import (
	"context"
	"sync"
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

	closed    bool
	closelock sync.Mutex
}

func (v *FabricRuntime) GetTimeout() time.Duration {
	return v.defaultTimeout
}

func (v *FabricRuntime) SetDefaultTimeout(t time.Duration) {
	v.defaultTimeout = t
}

func (v *FabricRuntime) Close() error {
	v.closelock.Lock()
	defer v.closelock.Unlock()
	if v.closed {
		return nil
	}
	v.closed = true

	v.hub.Close()
	return nil
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

type ServiceContext struct {
}

type StatelessUserServiceInstance interface {
	Run(ctx context.Context) error
	Close(ctx context.Context) error
	Abort() error
}

func (f *comFabricStatelessServiceFactoryGoProxy) init() {

}

func (v *comFabricStatelessServiceFactoryGoProxy) CreateInstance(
	_ *ole.IUnknown,
	serviceTypeName *uint16,
	serviceName *uint16,
	initializationDataLength uint32,
	initializationData *byte,
	partitionId windows.GUID,
	instanceId int64,
	serviceInstance **comFabricStatelessServiceInstance,
) uintptr {
	return 0
}

func (v *comFabricStatelessServiceInstanceGoProxy) init() {
}

func (v *comFabricStatelessServiceInstanceGoProxy) BeginOpen(
	_ *ole.IUnknown,
	partition *comFabricStatelessServicePartition,
	callback *comFabricAsyncOperationCallback,
	context **comFabricAsyncOperationContext,
) uintptr {
	return 0
}
func (v *comFabricStatelessServiceInstanceGoProxy) EndOpen(
	_ *ole.IUnknown,
	context *comFabricAsyncOperationContext,
	serviceAddress **comFabricStringResult,
) uintptr {
	return 0
}
func (v *comFabricStatelessServiceInstanceGoProxy) BeginClose(
	_ *ole.IUnknown,
	callback *comFabricAsyncOperationCallback,
	context **comFabricAsyncOperationContext,
) uintptr {
	return 0
}
func (v *comFabricStatelessServiceInstanceGoProxy) EndClose(
	_ *ole.IUnknown,
	context *comFabricAsyncOperationContext,
) uintptr {
	return 0
}
func (v *comFabricStatelessServiceInstanceGoProxy) Abort(
	_ *ole.IUnknown,
) uintptr {
	return 0
}
