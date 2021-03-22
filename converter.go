package fabric

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

func toUnsafePointer(s interface{}) unsafe.Pointer {
	if s == nil {
		return nil
	}
	switch s.(type) {
	case string:
		t, _ := windows.UTF16PtrFromString(s.(string))
		return unsafe.Pointer(t)
	case *FabricX509Credentials:
		t := s.(*FabricX509Credentials)
		return unsafe.Pointer(t.toInnerStruct())
	case FabricX509Credentials:
		t := s.(FabricX509Credentials)
		return unsafe.Pointer(t.toInnerStruct())
	}
	panic("not impl")
}

func fromUnsafePointer(unsafe.Pointer) interface{} {
	panic("not impl")
}
