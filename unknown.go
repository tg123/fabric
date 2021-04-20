package fabric

import (
	"sync"
	"unsafe"

	ole "github.com/go-ole/go-ole"
)

var (
	cgorefLock sync.Mutex
	cgoref     = map[unsafe.Pointer]interface{}{}
)

type goIUnknown struct {
	iid     string // current support 1 iid only
	ref     int32
	reflock sync.Mutex
}

func attachIUnknown(iid string, vtbl *ole.IUnknownVtbl) *goIUnknown {
	un := &goIUnknown{iid: iid}
	vtbl.QueryInterface = createCallback(un.queryInterface)
	vtbl.AddRef = createCallback(un.addRef)
	vtbl.Release = createCallback(un.release)
	return un
}

func (v *goIUnknown) queryInterface(this *ole.IUnknown, iid *ole.GUID, punk **ole.IUnknown) uintptr {
	*punk = nil
	if ole.IsEqualGUID(iid, ole.IID_IUnknown) {
		v.addRef(this)
		*punk = this
		return ole.S_OK
	}
	s, _ := ole.StringFromCLSID(iid)
	if s == v.iid {
		v.addRef(this)
		*punk = this
		return ole.S_OK
	}
	return ole.E_NOINTERFACE
}

func (v *goIUnknown) addRef(this *ole.IUnknown) uintptr {
	p := unsafe.Pointer(this)
	cgorefLock.Lock()
	cgoref[p] = this
	cgorefLock.Unlock()
	v.reflock.Lock()
	defer v.reflock.Unlock()
	v.ref++
	return uintptr(v.ref)
}

func (v *goIUnknown) release(this *ole.IUnknown) uintptr {
	p := unsafe.Pointer(this)
	v.reflock.Lock()
	defer v.reflock.Unlock()
	v.ref--
	if v.ref <= 0 {
		cgorefLock.Lock()
		delete(cgoref, p)
		cgorefLock.Unlock()
	}
	return uintptr(v.ref)
}
