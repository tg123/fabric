package fabric

import (
	"time"
	"unicode/utf16"
	"unsafe"
)

func toUnsafePointer(s interface{}) unsafe.Pointer {
	if s == nil {
		return nil
	}
	switch s.(type) {
	case string:
		t := utf16PtrFromString(s.(string))
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

func fromUnsafePointer(pointer unsafe.Pointer) interface{} {
	panic("not impl")
}

// from https://github.com/golang/sys/blob/1e4c9ba3b0c4/windows/types_windows.go#L759 to
type filetime struct {
	LowDateTime  uint32
	HighDateTime uint32
}

// Nanoseconds returns Filetime ft in nanoseconds
// since Epoch (00:00:00 UTC, January 1, 1970).
func (ft *filetime) Nanoseconds() int64 {
	// 100-nanosecond intervals since January 1, 1601
	nsec := int64(ft.HighDateTime)<<32 + int64(ft.LowDateTime)
	// change starting time to the Epoch (00:00:00 UTC, January 1, 1970)
	nsec -= 116444736000000000
	// convert into nanoseconds
	nsec *= 100
	return nsec
}

func (ft *filetime) ToTime() time.Time {
	return time.Unix(0, ft.Nanoseconds())
}

func timeToFiletime(t time.Time) filetime {
	return nsecToFiletime(t.UnixNano())
}

func nsecToFiletime(nsec int64) (ft filetime) {
	// convert into 100-nanosecond
	nsec /= 100
	// change starting time to January 1, 1601
	nsec += 116444736000000000
	// split into high / low
	ft.LowDateTime = uint32(nsec & 0xffffffff)
	ft.HighDateTime = uint32(nsec >> 32 & 0xffffffff)
	return ft
}

// https://github.com/golang/sys/blob/1e4c9ba3b0c4fcddbe90893331bdc829813066a1/windows/syscall_windows.go#L88
func utf16PtrFromString(s string) *uint16 {
	for i := 0; i < len(s); i++ {
		if s[i] == 0 {
			return nil
		}
	}
	a := utf16.Encode([]rune(s + "\x00"))
	return &a[0]
}

func utf16PtrToString(p *uint16) string {
	if p == nil {
		return ""
	}
	if *p == 0 {
		return ""
	}

	// Find NUL terminator.
	n := 0
	for ptr := unsafe.Pointer(p); *(*uint16)(ptr) != 0; n++ {
		ptr = unsafe.Pointer(uintptr(ptr) + unsafe.Sizeof(*p))
	}

	return string(utf16.Decode(unsafe.Slice(p, n)))
}
