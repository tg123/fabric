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
	vtbl     *comIFabricAsyncOperationCallbackVtbl
	ref      int32
	callback func(ctx *comIFabricAsyncOperationContext)
}

type comIFabricAsyncOperationCallbackVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
	Invoke         uintptr
}

func newFabricAsyncOperationCallback(fn func(ctx *comIFabricAsyncOperationContext)) *comIFabricAsyncOperationCallback {
	cb := &comIFabricAsyncOperationCallback{}
	cb.vtbl = &comIFabricAsyncOperationCallbackVtbl{}
	cb.vtbl.QueryInterface = syscall.NewCallback(cb.queryInterface)
	cb.vtbl.AddRef = syscall.NewCallback(cb.addRef)
	cb.vtbl.Release = syscall.NewCallback(cb.release)
	cb.vtbl.Invoke = syscall.NewCallback(cb.invoke)
	cb.callback = fn
	return cb
}

func (v *comIFabricAsyncOperationCallback) queryInterface(this *ole.IUnknown, iid *ole.GUID, punk **ole.IUnknown) uintptr {
	s, _ := ole.StringFromCLSID(iid)
	*punk = nil
	if ole.IsEqualGUID(iid, ole.IID_IUnknown) {
		v.addRef(this)
		*punk = this
		return ole.S_OK
	}
	if s == "{86F08D7E-14DD-4575-8489-B1D5D679029C}" {
		v.addRef(this)
		*punk = this
		return ole.S_OK
	}
	return ole.E_NOINTERFACE
}

func (v *comIFabricAsyncOperationCallback) addRef(this *ole.IUnknown) uintptr {
	pthis := (*comIFabricAsyncOperationCallback)(unsafe.Pointer(this))
	pthis.ref++
	return uintptr(pthis.ref)
}

func (v *comIFabricAsyncOperationCallback) release(this *ole.IUnknown) uintptr {
	pthis := (*comIFabricAsyncOperationCallback)(unsafe.Pointer(this))
	pthis.ref--
	return uintptr(pthis.ref)
}

func (v *comIFabricAsyncOperationCallback) invoke(this *ole.IUnknown, ctx *ole.IUnknown) uintptr {
	pthis := (*comIFabricAsyncOperationCallback)(unsafe.Pointer(this))

	ctx.AddRef()
	defer ctx.Release()
	pthis.callback((*comIFabricAsyncOperationContext)(unsafe.Pointer(ctx)))

	return ole.S_OK
}
