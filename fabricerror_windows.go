// +build windows

package fabric

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	fabricCommonDll               = windows.NewLazyDLL("FabricCommon.dll")
	fabricGetLastErrorMessageProc = fabricCommonDll.NewProc("FabricGetLastErrorMessage")
)

var errComNotImpl = fmt.Errorf("operation not supported on this fabric version")

func fabricGetLastError() string {
	var result *comFabricStringResult
	hr, _, _ := fabricGetLastErrorMessageProc.Call(uintptr(unsafe.Pointer(&result)))

	if hr != 0 {
		return ""
	}

	msg, _ := result.GetString()
	return msg
}
