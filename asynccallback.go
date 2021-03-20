package fabric

import (
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
