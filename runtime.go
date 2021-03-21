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
	r, _, err := fabricFabricCreateRuntimeProc.Call(uintptr(unsafe.Pointer(clzid)), uintptr(unsafe.Pointer(&com)))

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
	ServiceTypeName    string
	ServiceName        string
	InitializationData []byte
	PartitionId        windows.GUID
	InstanceId         int64
}

// TODO genrate interface from mkidl
type FabricStatelessServicePartition interface {
	GetPartitionInfo() (bufferedValue *FabricServicePartitionInformation, err error)
}

type StatelessServiceInstance interface {
	Open(ctx context.Context, partition FabricStatelessServicePartition) (string, error)
	Close(ctx context.Context) error
	Abort() error
}

func (v *FabricRuntime) RegisterStatelessService(serviceTypeName string, builder func(ServiceContext) (StatelessServiceInstance, error)) error {
	b := newComFabricStatelessServiceFactory(builder)
	return v.registerStatelessServiceFactory(serviceTypeName, b)
}

func (f *comFabricStatelessServiceFactoryGoProxy) init() {
}

func (v *comFabricStatelessServiceFactoryGoProxy) CreateInstance(
	_ *ole.IUnknown,
	serviceTypeName *uint16,
	serviceName *uint16,
	initializationDataLength uint32,
	initializationData *byte,
	partitionId *windows.GUID,
	instanceId int64,
	serviceInstance **comFabricStatelessServiceInstance,
) uintptr {
	ctx := ServiceContext{
		ServiceTypeName: windows.UTF16PtrToString(serviceTypeName),
		ServiceName:     windows.UTF16PtrToString(serviceName),
		PartitionId:     *partitionId,
		InstanceId:      instanceId,
	}
	sliceCast(unsafe.Pointer(&ctx.InitializationData), unsafe.Pointer(initializationData), int(initializationDataLength))

	inst, err := v.builder(ctx)

	if err != nil {
		return errorToHResult(err)
	}

	*serviceInstance = newComFabricStatelessServiceInstance(inst)
	return ole.S_OK
}

func (v *comFabricStatelessServiceInstanceGoProxy) init() {
}

func (v *comFabricStatelessServiceInstanceGoProxy) BeginOpen(
	_ *ole.IUnknown,
	partition *comFabricStatelessServicePartition,
	callback *comFabricAsyncOperationCallback,
	asyncContext **comFabricAsyncOperationContext,
) uintptr {
	partition.AddRef()
	return asyncRun(func(goctx context.Context) (interface{}, error) {
		defer partition.Release()
		return v.instance.Open(goctx, partition)
	}, callback, asyncContext)
}

func (v *comFabricStatelessServiceInstanceGoProxy) EndOpen(
	_ *ole.IUnknown,
	context *comFabricAsyncOperationContext,
	serviceAddress **comFabricStringResult,
) uintptr {
	<-context.proxy.goctx.Done()

	if context.proxy.resultHResult != ole.S_OK {
		return context.proxy.resultHResult
	}

	result, ok := context.proxy.result.(string)

	if !ok {
		return ole.E_UNEXPECTED
	}

	*serviceAddress = newComFabricStringResult(result)
	return ole.S_OK
}

func (v *comFabricStatelessServiceInstanceGoProxy) BeginClose(
	_ *ole.IUnknown,
	callback *comFabricAsyncOperationCallback,
	asyncContext **comFabricAsyncOperationContext,
) uintptr {
	return asyncRun(func(goctx context.Context) (interface{}, error) {
		return nil, v.instance.Close(goctx)
	}, callback, asyncContext)
}

func (v *comFabricStatelessServiceInstanceGoProxy) EndClose(
	_ *ole.IUnknown,
	context *comFabricAsyncOperationContext,
) uintptr {
	<-context.proxy.goctx.Done()
	return context.proxy.resultHResult
}

func (v *comFabricStatelessServiceInstanceGoProxy) Abort(
	_ *ole.IUnknown,
) uintptr {
	return errorToHResult(v.instance.Abort())
}
