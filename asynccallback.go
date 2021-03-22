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

func (a *comFabricAsyncOperationContextGoProxy) isCompleted() bool {
	return a.goctx.Err() != nil
}

func (a *comFabricAsyncOperationContextGoProxy) completeWith(result interface{}, hr uintptr) {
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.isCompleted() {
		return
	}

	a.result = result
	a.resultHResult = hr
	a.cancel()
}

func (a *comFabricAsyncOperationContextGoProxy) await() (result interface{}, hr uintptr) {
	<-a.goctx.Done()
	return a.result, a.resultHResult
}

func (a *comFabricAsyncOperationContextGoProxy) IsCompleted(_ *ole.IUnknown) uintptr {
	if a.isCompleted() {
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
	a.completeWith(nil, ole.E_ABORT)
	return ole.S_OK
}

func asyncRun(
	action func(goctx context.Context) (interface{}, error),
	callback *comFabricAsyncOperationCallback,
	asyncContext **comFabricAsyncOperationContext,
) uintptr {
	callback.AddRef()
	goctx, cancel := context.WithCancel(context.Background())
	ctx := newComFabricAsyncOperationContext(callback, goctx, cancel)
	go func() {
		defer callback.Release()
		go func() {
			r, err := action(goctx)
			ctx.proxy.completeWith(r, errorToHResult(err))
		}()

		<-goctx.Done()
		callback.Invoke(ctx)
	}()

	*asyncContext = ctx
	return ole.S_OK
}
