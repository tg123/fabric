package fabric

import (
	"context"
	"fmt"
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

func queryObject(com *ole.IUnknown, iid string, outptr unsafe.Pointer) error {
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

func toTimeout(ctx context.Context) time.Duration {
	deadline, ok := ctx.Deadline()
	if ok {
		return deadline.Sub(time.Now())
	}

	// TODO move to sf client var
	return 15 * time.Minute
}
