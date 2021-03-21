package fabric

import (
	"context"
	"fmt"
	"reflect"
	"time"
	"unsafe"

	ole "github.com/go-ole/go-ole"
	"github.com/pkg/errors"
	"golang.org/x/sys/windows"
)

var (
	fabricCommonDll               = windows.NewLazyDLL("FabricCommon.dll")
	fabricGetLastErrorMessageProc = fabricCommonDll.NewProc("FabricGetLastErrorMessage")
)

var errComNotImpl = fmt.Errorf("operation not supported on this fabric version")

type comCreator func(iid string, outptr unsafe.Pointer) error

func queryComObject(com *ole.IUnknown, iid string, outptr unsafe.Pointer) error {
	clzid, err := ole.IIDFromString(iid)
	if err != nil {
		return err
	}

	c, err := com.QueryInterface(clzid)
	if err != nil {
		return err
	}

	*(**ole.IUnknown)(outptr) = &c.IUnknown

	return nil
}

// not thread safe lock before use
func releaseComObject(com *ole.IUnknown) error {

	// TODO
	// according to the doc https://docs.microsoft.com/en-us/windows/win32/api/unknwn/nf-unknwn-iunknown-release
	// the return should be only for test
	// here we assume all our objects will ref == 0 eventually
	// however
	// now some com objects exposed and calling release is not guaranteed

	com.Release()
	return nil
}

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

	return fmt.Sprintf("error [%v] [0x%x]", c.String(), uint64(c))
}

func errno(hr uintptr, syserr error) error {
	if hr == 0 {
		return nil
	}

	return errors.Wrap(FabricErrorCode(hr), fabricGetLastError())
}

func waitch(ctx context.Context, ch <-chan error, sfctx *comFabricAsyncOperationContext, timeout time.Duration) (err error) {
	defer releaseComObject(&sfctx.IUnknown)
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

type withTimeout interface {
	GetTimeout() time.Duration
	SetDefaultTimeout(time.Duration)
}

func toTimeout(ctx context.Context, v withTimeout) time.Duration {
	deadline, ok := ctx.Deadline()
	if ok {
		return deadline.Sub(time.Now())
	}

	return v.GetTimeout()
}

func sliceCast(dst, src unsafe.Pointer, len int) {
	slice := (*reflect.SliceHeader)(dst)
	slice.Data = uintptr(src)
	slice.Len = len
	slice.Cap = len
}

func (v *comFabricStringResultGoProxy) init() {
}

func (v *comFabricStringResultGoProxy) GetString(_ *ole.IUnknown) uintptr {
	s, _ := windows.UTF16PtrFromString(v.result)
	return uintptr(unsafe.Pointer(s))
}

func errorToHResult(err error) uintptr {
	if err != nil {
		if err1, ok := err.(FabricErrorCode); ok {
			return uintptr(err1)
		}

		return ole.E_FAIL
	}

	return ole.S_OK
}
