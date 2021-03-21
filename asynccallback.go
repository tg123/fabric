package fabric

import (
	"context"
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
	a.resultHResult = ole.E_ABORT
	a.result = nil
	a.cancel()
	return ole.S_OK
}

func asyncRun(
	action func(goctx context.Context) (interface{}, error),
	callback *comFabricAsyncOperationCallback,
	asyncContext **comFabricAsyncOperationContext,
) uintptr {
	callback.AddRef()
	goctx, cancel := context.WithCancel(context.Background())
	ctx := newComFabricAsyncOperationContext(callback, nil, ole.S_OK, goctx, cancel)
	go func() {
		defer callback.Release()
		go func() {
			r, err := action(goctx)
			ctx.proxy.result = r
			ctx.proxy.resultHResult = errorToHResult(err)
			cancel()
		}()

		<-goctx.Done()
		callback.Invoke(ctx)
	}()

	*asyncContext = ctx
	return ole.S_OK
}
