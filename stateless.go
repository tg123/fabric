package fabric

import (
	"context"
	"unsafe"

	ole "github.com/go-ole/go-ole"
	"golang.org/x/sys/windows"
)

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
	r, hr := context.proxy.await()
	<-context.proxy.goctx.Done()

	if hr != ole.S_OK {
		return hr
	}

	result, ok := r.(string)

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
	_, hr := context.proxy.await()
	return hr
}

func (v *comFabricStatelessServiceInstanceGoProxy) Abort(
	_ *ole.IUnknown,
) uintptr {
	return errorToHResult(v.instance.Abort())
}
