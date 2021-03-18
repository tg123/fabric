package fabric

import (
	"syscall"

	ole "github.com/go-ole/go-ole"
)

// TODO generate go side ?
type comIFabricAsyncOperationCallback struct {
	vtbl       *comIFabricAsyncOperationCallbackVtbl
	callback   func(ctx *comFabricAsyncOperationContext)
	unknownref *goIUnknown
}

type comIFabricAsyncOperationCallbackVtbl struct {
	goIUnknownVtbl
	Invoke uintptr
}

func newFabricAsyncOperationCallback(fn func(ctx *comFabricAsyncOperationContext)) *comIFabricAsyncOperationCallback {
	cb := &comIFabricAsyncOperationCallback{}
	cb.vtbl = &comIFabricAsyncOperationCallbackVtbl{}
	cb.unknownref = attachIUnknown("{86F08D7E-14DD-4575-8489-B1D5D679029C}", &cb.vtbl.goIUnknownVtbl)
	cb.vtbl.Invoke = syscall.NewCallback(cb.invoke)
	cb.callback = fn
	return cb
}

func (v *comIFabricAsyncOperationCallback) invoke(this *ole.IUnknown, ctx *comFabricAsyncOperationContext) uintptr {
	if ctx == nil {
		return ole.E_POINTER
	}
	ctx.AddRef()
	defer ctx.Release()
	v.callback(ctx)
	return ole.S_OK
}
