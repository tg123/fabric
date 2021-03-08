package fabric

import (
	"context"
	"fmt"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	fabricCommonDll               = windows.MustLoadDLL("FabricCommon.dll")
	fabricGetLastErrorMessageProc = fabricCommonDll.MustFindProc("FabricGetLastErrorMessage")
)

func fabricGetLastError() string {
	var result *comFabricStringResult
	hr, _, _ := fabricGetLastErrorMessageProc.Call(uintptr(unsafe.Pointer(&result)))

	if hr != 0 {
		return ""
	}

	msg, _ := result.GetString()
	return msg
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

func waitch(ctx context.Context, ch <-chan error, sfctx *comIFabricAsyncOperationContext, timeout time.Duration) (err error) {
	select {
	case err = <-ch:
		return
	case <-ctx.Done():
		sfctx.Cancel()
		err = ctx.Err()
		return
	case <-time.After(timeout):
		sfctx.Cancel()
		err = FabricErrorTimeout
		return
	}
}
