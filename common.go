package fabric

import (
	"fmt"
	"syscall"
	"unsafe"

	ole "github.com/go-ole/go-ole"
	"golang.org/x/sys/windows"
)

var (
	fabricCommonDll               = windows.MustLoadDLL("FabricCommon.dll")
	fabricGetLastErrorMessageProc = fabricCommonDll.MustFindProc("FabricGetLastErrorMessage")
)

type comIFabricStringResult struct {
	ole.IUnknown
}

type comIFabricStringResultVtbl struct {
	ole.IUnknownVtbl
	GetString uintptr
}

func (v *comIFabricStringResult) VTable() *comIFabricStringResultVtbl {
	return (*comIFabricStringResultVtbl)(unsafe.Pointer(v.RawVTable))
}

func fabricGetLastError() string {
	var result *comIFabricStringResult
	hr, _, _ := fabricGetLastErrorMessageProc.Call(uintptr(unsafe.Pointer(&result)))

	// ..... wtf
	if hr != 0 {
		return ""
	}

	hr, _, _ = syscall.Syscall(
		uintptr(result.VTable().GetString),
		1,
		uintptr(unsafe.Pointer(result)),
		0,
		0)

	if hr == 0 {
		return ""
	}

	return windows.UTF16PtrToString((*uint16)(unsafe.Pointer(hr)))
}

func (c FabricErrorCode) Error() string {
	if c == 0 {
		return ""
	}

	return fmt.Sprintf("error [%v] [0x%x] msg: [%v]", c.String(), uint64(c), fabricGetLastError())
}

func errno(hr uintptr, syserr error) error {
	if hr == 0 {
		return nil
	}

	return FabricErrorCode(hr)
}
