package fabric

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// TODO clean up, those struct copied from generated files

type fabricX509Credentials struct {
	AllowedCommonNames []string
	FindType           FabricX509FindType
	FindValue          interface{}
	StoreLocation      FabricX509StoreLocation
	StoreName          string
	ProtectionLevel    FabricProtectionLevel
}

func (obj *fabricX509Credentials) toInnerStruct() *innerFabricX509Credentials {
	if obj == nil {
		return nil
	}
	dst := innerFabricX509Credentials{}
	var lst_1 []*uint16
	for _, item := range obj.AllowedCommonNames {
		s, _ := windows.UTF16PtrFromString(item)
		lst_1 = append(lst_1, s)
	}
	if len(lst_1) > 0 {
		dst.AllowedCommonNames = &lst_1[0]
	}
	dst.AllowedCommonNameCount = uint32(len(lst_1))
	dst.FindType = obj.FindType
	dst.FindValue = unsafe.Pointer(&obj.FindValue)
	{
		tmp_3, ok_3 := (obj.FindValue).(string)
		if ok_3 {
			tmps_3, _ := windows.UTF16PtrFromString(tmp_3)
			dst.FindValue = unsafe.Pointer(tmps_3)
		}
	}

	dst.StoreLocation = obj.StoreLocation
	tmp_5, _ := windows.UTF16PtrFromString(obj.StoreName)
	dst.StoreName = tmp_5
	dst.ProtectionLevel = obj.ProtectionLevel
	return &dst
}

type fabricX509CredentialsEx1 struct {
	IssuerThumbprints []string
}

func (obj *fabricX509CredentialsEx1) toInnerStruct() *innerFabricX509CredentialsEx1 {
	if obj == nil {
		return nil
	}
	dst := innerFabricX509CredentialsEx1{}
	var lst_1 []*uint16
	for _, item := range obj.IssuerThumbprints {
		s, _ := windows.UTF16PtrFromString(item)
		lst_1 = append(lst_1, s)
	}
	if len(lst_1) > 0 {
		dst.IssuerThumbprints = &lst_1[0]
	}
	dst.IssuerThumbprintCount = uint32(len(lst_1))
	return &dst
}

type fabricX509CredentialsEx2 struct {
	RemoteCertThumbprints []string
	RemoteX509NameCount   uint32
	RemoteX509Names       *innerFabricX509Name
	FindValueSecondary    interface{}
}

func (obj *fabricX509CredentialsEx2) toInnerStruct() *innerFabricX509CredentialsEx2 {
	if obj == nil {
		return nil
	}
	dst := innerFabricX509CredentialsEx2{}
	var lst_1 []*uint16
	for _, item := range obj.RemoteCertThumbprints {
		s, _ := windows.UTF16PtrFromString(item)
		lst_1 = append(lst_1, s)
	}
	if len(lst_1) > 0 {
		dst.RemoteCertThumbprints = &lst_1[0]
	}
	dst.RemoteCertThumbprintCount = uint32(len(lst_1))
	dst.RemoteX509NameCount = 0
	// dst.RemoteX509Names
	dst.FindValueSecondary = unsafe.Pointer(&obj.FindValueSecondary)
	{
		tmp_4, ok_4 := (obj.FindValueSecondary).(string)
		if ok_4 {
			tmps_4, _ := windows.UTF16PtrFromString(tmp_4)
			dst.FindValueSecondary = unsafe.Pointer(tmps_4)
		}
	}
	return &dst
}

type fabricSecurityCredentials struct {
	Kind  FabricSecurityCredentialKind
	Value unsafe.Pointer
}

// hack due to interface{} to unsafe
func (v *comFabricClientSettings) innerSetSecurityCredentials(
	cred *X509Credentials,
) (err error) {

	c := fabricX509Credentials{
		FindType:        cred.FindType,
		FindValue:       cred.FindValue,
		StoreLocation:   cred.StoreLocation,
		StoreName:       cred.StoreName,
		ProtectionLevel: cred.ProtectionLevel,
	}

	innerc := c.toInnerStruct()

	ex1 := fabricX509CredentialsEx1{
		IssuerThumbprints: cred.IssuerThumbprints,
	}

	innerex1 := ex1.toInnerStruct()
	innerc.Reserved = unsafe.Pointer(innerex1)

	ex2 := fabricX509CredentialsEx2{
		RemoteCertThumbprints: cred.RemoteCertThumbprints,
	}

	innerex2 := ex2.toInnerStruct()
	innerex2.FindValueSecondary = nil

	innerex1.Reserved = unsafe.Pointer(innerex2)

	securityCredentials := fabricSecurityCredentials{
		Kind:  FabricSecurityCredentialKindX509,
		Value: unsafe.Pointer(innerc),
	}

	hr, _, err1 := syscall.Syscall(
		v.vtable().SetSecurityCredentials,
		2,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(&securityCredentials)),
		0,
	)
	if hr != 0 {
		err = errno(hr, err1)
		return
	}
	return
}

type X509Credentials struct {
	// AllowedCommonNames []string
	FindType              FabricX509FindType
	FindValue             interface{}
	StoreLocation         FabricX509StoreLocation
	StoreName             string
	ProtectionLevel       FabricProtectionLevel
	RemoteCertThumbprints []string
	IssuerThumbprints     []string
}
