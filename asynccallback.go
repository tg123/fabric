package fabric

import (
	"syscall"
	"unsafe"

	ole "github.com/go-ole/go-ole"
)

type comIFabricAsyncOperationContext struct {
	ole.IUnknown
}

type comIFabricAsyncOperationContextVtbl struct {
	ole.IUnknownVtbl
	IsCompleted            uintptr
	CompletedSynchronously uintptr
	GetCallback            uintptr
	Cancel                 uintptr
}

func (v *comIFabricAsyncOperationContext) VTable() *comIFabricAsyncOperationContextVtbl {
	return (*comIFabricAsyncOperationContextVtbl)(unsafe.Pointer(v.RawVTable))
}

func (v *comIFabricAsyncOperationContext) IsCompleted() bool {
	hr, _, _ := syscall.Syscall(
		v.VTable().IsCompleted,
		1,
		uintptr(unsafe.Pointer(v)),
		0,
		0)

	return hr != 0
}

func (v *comIFabricAsyncOperationContext) CompletedSynchronously() bool {
	hr, _, _ := syscall.Syscall(
		v.VTable().CompletedSynchronously,
		1,
		uintptr(unsafe.Pointer(v)),
		0,
		0)

	return hr != 0
}

func (v *comIFabricAsyncOperationContext) Cancel() (err error) {
	hr, _, err1 := syscall.Syscall(
		v.VTable().Cancel,
		1,
		uintptr(unsafe.Pointer(v)),
		0,
		0,
	)
	if hr != 0 {
		err = errno(hr, err1)
		return
	}
	return
}

type comIFabricAsyncOperationCallback struct {
	vtbl       *comIFabricAsyncOperationCallbackVtbl
	callback   func(ctx *comIFabricAsyncOperationContext)
	unknownref *goIUnknown
}

type comIFabricAsyncOperationCallbackVtbl struct {
	goIUnknownVtbl
	Invoke uintptr
}

func newFabricAsyncOperationCallback(fn func(ctx *comIFabricAsyncOperationContext)) *comIFabricAsyncOperationCallback {
	cb := &comIFabricAsyncOperationCallback{}
	cb.vtbl = &comIFabricAsyncOperationCallbackVtbl{}
	cb.unknownref = attachIUnknown("{86F08D7E-14DD-4575-8489-B1D5D679029C}", &cb.vtbl.goIUnknownVtbl)
	cb.vtbl.Invoke = syscall.NewCallback(cb.invoke)
	cb.callback = fn
	return cb
}

func (v *comIFabricAsyncOperationCallback) invoke(this *ole.IUnknown, ctx *comIFabricAsyncOperationContext) uintptr {
	if ctx == nil {
		return ole.E_POINTER
	}
	ctx.AddRef()
	defer ctx.Release()
	v.callback(ctx)
	return ole.S_OK
}
