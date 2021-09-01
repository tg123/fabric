package fabric

import (
	"context"
	"io"
	"sync"
	"time"
	"unsafe"

	ole "github.com/go-ole/go-ole"
)

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

func (v *comFabricStringResultGoProxy) init() {
}

func (v *comFabricStringResultGoProxy) GetString(_ *ole.IUnknown) uintptr {
	return uintptr(unsafe.Pointer(utf16PtrFromString(v.result)))
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

type coclz struct {
	hub            io.Closer
	defaultTimeout time.Duration

	closed    bool
	closelock sync.Mutex

	deferclose []func()
}

func (v *coclz) GetTimeout() time.Duration {
	return v.defaultTimeout
}

func (v *coclz) SetDefaultTimeout(t time.Duration) {
	v.defaultTimeout = t
}

func (v *coclz) Close() error {
	v.closelock.Lock()
	defer v.closelock.Unlock()
	if v.closed {
		return nil
	}
	v.closed = true

	v.hub.Close()
	for _, cf := range v.deferclose {
		cf()
	}
	return nil
}
