package fabric

import (
	"unsafe"

	ole "github.com/go-ole/go-ole"
)

func (cb *comFabricAsyncOperationCallbackGoProxy) init() {

}

func (cb *comFabricAsyncOperationCallbackGoProxy) Invoke(this *ole.IUnknown, ctx *comFabricAsyncOperationContext) uintptr {
	if ctx == nil {
		return ole.E_POINTER
	}
	ctx.AddRef()
	defer ctx.Release()
	cb.callback(ctx)
	return ole.S_OK
}

func (a *comFabricAsyncOperationContextGoProxy) init() {

}

func (a *comFabricAsyncOperationContextGoProxy) IsCompleted(_ *ole.IUnknown) uintptr {
	done := a.goctx.Err() != nil
	if done {
		return 1 // true
	}
	return 0 // false
}

func (a *comFabricAsyncOperationContextGoProxy) CompletedSynchronously(_ *ole.IUnknown) uintptr {
	return 0 // false
}

func (a *comFabricAsyncOperationContextGoProxy) GetCallback(_ *ole.IUnknown, callback **comFabricAsyncOperationCallback) uintptr {
	return uintptr(unsafe.Pointer(a.nativeCallback))
}

func (a *comFabricAsyncOperationContextGoProxy) Cancel(_ *ole.IUnknown) uintptr {
	a.cancel()
	return ole.S_OK
}
