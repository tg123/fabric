// Code generated by "go run github.com/tg123/fabric/mkidl"; DO NOT EDIT.
package fabric

import (
	"context"
	"github.com/go-ole/go-ole"
	"golang.org/x/sys/windows"
	"sync"
	"syscall"
	"unsafe"
)

type comFabricAsyncOperationCallbackGoProxy struct {
	unknownref *goIUnknown
	callback   func(ctx *comFabricAsyncOperationContext)
}

func newComFabricAsyncOperationCallback(
	callback func(ctx *comFabricAsyncOperationContext),
) *comFabricAsyncOperationCallback {
	com := &comFabricAsyncOperationCallback{}
	*(**comFabricAsyncOperationCallbackVtbl)(unsafe.Pointer(com)) = &comFabricAsyncOperationCallbackVtbl{}
	vtbl := com.vtable()
	com.proxy.unknownref = attachIUnknown("{86F08D7E-14DD-4575-8489-B1D5D679029C}", &vtbl.IUnknownVtbl)
	vtbl.Invoke = syscall.NewCallback(com.proxy.Invoke)

	com.proxy.callback = callback

	com.proxy.init()
	return com
}

/*
func (v *comFabricAsyncOperationCallbackGoProxy) Invoke(
_ *ole.IUnknown,
context *comFabricAsyncOperationContext,
) uintptr { return 0}
*/

type comFabricAsyncOperationCallback struct {
	ole.IUnknown
	proxy comFabricAsyncOperationCallbackGoProxy
}

type comFabricAsyncOperationCallbackVtbl struct {
	ole.IUnknownVtbl
	Invoke uintptr
}

func (v *comFabricAsyncOperationCallback) vtable() *comFabricAsyncOperationCallbackVtbl {
	return (*comFabricAsyncOperationCallbackVtbl)(unsafe.Pointer(v.RawVTable))
}

func (v *comFabricAsyncOperationCallback) Invoke(
	context *comFabricAsyncOperationContext,
) (rt interface{}, err error) {
	hr, _, err1 := syscall.Syscall(
		v.vtable().Invoke,
		2,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(context)),
		0,
	)
	if hr == 0 {
		err = err1
		return
	}

	tmp := (unsafe.Pointer)(unsafe.Pointer(hr))

	rt = fromUnsafePointer(tmp)
	return
}

type comFabricAsyncOperationContextGoProxy struct {
	unknownref     *goIUnknown
	nativeCallback *comFabricAsyncOperationCallback
	result         interface{}
	resultHResult  uintptr
	lock           sync.Mutex
	goctx          context.Context
	cancel         context.CancelFunc
}

func newComFabricAsyncOperationContext(
	nativeCallback *comFabricAsyncOperationCallback,

	goctx context.Context,
	cancel context.CancelFunc,
) *comFabricAsyncOperationContext {
	com := &comFabricAsyncOperationContext{}
	*(**comFabricAsyncOperationContextVtbl)(unsafe.Pointer(com)) = &comFabricAsyncOperationContextVtbl{}
	vtbl := com.vtable()
	com.proxy.unknownref = attachIUnknown("{841720BF-C9E8-4E6F-9C3F-6B7F4AC73BCD}", &vtbl.IUnknownVtbl)
	vtbl.IsCompleted = syscall.NewCallback(com.proxy.IsCompleted)
	vtbl.CompletedSynchronously = syscall.NewCallback(com.proxy.CompletedSynchronously)
	vtbl.get_Callback = syscall.NewCallback(com.proxy.GetCallback)
	vtbl.Cancel = syscall.NewCallback(com.proxy.Cancel)

	com.proxy.nativeCallback = nativeCallback

	com.proxy.goctx = goctx
	com.proxy.cancel = cancel

	com.proxy.init()
	return com
}

/*
func (v *comFabricAsyncOperationContextGoProxy) IsCompleted(
_ *ole.IUnknown,
) uintptr { return 0}
func (v *comFabricAsyncOperationContextGoProxy) CompletedSynchronously(
_ *ole.IUnknown,
) uintptr { return 0}
func (v *comFabricAsyncOperationContextGoProxy) GetCallback(
_ *ole.IUnknown,
callback **comFabricAsyncOperationCallback,
) uintptr { return 0}
func (v *comFabricAsyncOperationContextGoProxy) Cancel(
_ *ole.IUnknown,
) uintptr { return 0}
*/

type comFabricAsyncOperationContext struct {
	ole.IUnknown
	proxy comFabricAsyncOperationContextGoProxy
}

type comFabricAsyncOperationContextVtbl struct {
	ole.IUnknownVtbl
	IsCompleted            uintptr
	CompletedSynchronously uintptr
	get_Callback           uintptr
	Cancel                 uintptr
}

func (v *comFabricAsyncOperationContext) vtable() *comFabricAsyncOperationContextVtbl {
	return (*comFabricAsyncOperationContextVtbl)(unsafe.Pointer(v.RawVTable))
}

func (v *comFabricAsyncOperationContext) IsCompleted() (rt bool, err error) {
	hr, _, err1 := syscall.Syscall(
		v.vtable().IsCompleted,
		1,
		uintptr(unsafe.Pointer(v)),
		0,
		0,
	)
	_ = err1
	rt = hr != 0
	return
}
func (v *comFabricAsyncOperationContext) CompletedSynchronously() (rt bool, err error) {
	hr, _, err1 := syscall.Syscall(
		v.vtable().CompletedSynchronously,
		1,
		uintptr(unsafe.Pointer(v)),
		0,
		0,
	)
	_ = err1
	rt = hr != 0
	return
}
func (v *comFabricAsyncOperationContext) GetCallback() (callback *comFabricAsyncOperationCallback, err error) {
	var p_0 *comFabricAsyncOperationCallback
	defer func() {
		callback = p_0
	}()
	hr, _, err1 := syscall.Syscall(
		v.vtable().get_Callback,
		2,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(&p_0)),
		0,
	)
	if hr != 0 {
		err = errno(hr, err1)
		return
	}
	return
}
func (v *comFabricAsyncOperationContext) Cancel() (err error) {
	hr, _, err1 := syscall.Syscall(
		v.vtable().Cancel,
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

type comFabricStringResultGoProxy struct {
	unknownref *goIUnknown
	result     string
}

func newComFabricStringResult(
	result string,
) *comFabricStringResult {
	com := &comFabricStringResult{}
	*(**comFabricStringResultVtbl)(unsafe.Pointer(com)) = &comFabricStringResultVtbl{}
	vtbl := com.vtable()
	com.proxy.unknownref = attachIUnknown("{4AE69614-7D0F-4CD4-B836-23017000D132}", &vtbl.IUnknownVtbl)
	vtbl.get_String = syscall.NewCallback(com.proxy.GetString)

	com.proxy.result = result

	com.proxy.init()
	return com
}

/*
func (v *comFabricStringResultGoProxy) GetString(
_ *ole.IUnknown,
) uintptr { return 0}
*/

type comFabricStringResult struct {
	ole.IUnknown
	proxy comFabricStringResultGoProxy
}

type comFabricStringResultVtbl struct {
	ole.IUnknownVtbl
	get_String uintptr
}

func (v *comFabricStringResult) vtable() *comFabricStringResultVtbl {
	return (*comFabricStringResultVtbl)(unsafe.Pointer(v.RawVTable))
}

func (v *comFabricStringResult) GetString() (rt string, err error) {
	hr, _, err1 := syscall.Syscall(
		v.vtable().get_String,
		1,
		uintptr(unsafe.Pointer(v)),
		0,
		0,
	)
	if hr == 0 {
		err = err1
		return
	}

	tmp := (*uint16)(unsafe.Pointer(hr))

	rt = windows.UTF16PtrToString(tmp)
	return
}

type comFabricStringListResult struct {
	ole.IUnknown
}

type comFabricStringListResultVtbl struct {
	ole.IUnknownVtbl
	GetStrings uintptr
}

func (v *comFabricStringListResult) vtable() *comFabricStringListResultVtbl {
	return (*comFabricStringListResultVtbl)(unsafe.Pointer(v.RawVTable))
}

type comFabricGetReplicatorStatusResult struct {
	ole.IUnknown
}

type comFabricGetReplicatorStatusResultVtbl struct {
	ole.IUnknownVtbl
	get_ReplicatorStatus uintptr
}

func (v *comFabricGetReplicatorStatusResult) vtable() *comFabricGetReplicatorStatusResultVtbl {
	return (*comFabricGetReplicatorStatusResultVtbl)(unsafe.Pointer(v.RawVTable))
}

func (v *comFabricGetReplicatorStatusResult) GetReplicatorStatus() (rt *FabricReplicatorStatusQueryResult, err error) {
	hr, _, err1 := syscall.Syscall(
		v.vtable().get_ReplicatorStatus,
		1,
		uintptr(unsafe.Pointer(v)),
		0,
		0,
	)
	if hr == 0 {
		err = err1
		return
	}

	tmp := (*innerFabricReplicatorStatusQueryResult)(unsafe.Pointer(hr))

	rt = tmp.toGoStruct()
	return
}
