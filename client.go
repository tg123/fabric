package fabric

import (
	"fmt"
	"unsafe"

	ole "github.com/go-ole/go-ole"
	"golang.org/x/sys/windows"
)

var (
	fabricClientdll             = windows.MustLoadDLL("FabricClient.dll")
	fabricCreateLocalClientProc = fabricClientdll.MustFindProc("FabricCreateLocalClient")
	fabricCreateClientProc      = fabricClientdll.MustFindProc("FabricCreateClient")
)

func createLocalClient(iid string, p unsafe.Pointer) error {
	clzid, err := ole.IIDFromString(iid)
	if err != nil {
		return err
	}

	r, _, err := fabricCreateLocalClientProc.Call(uintptr(unsafe.Pointer(clzid)), uintptr(p))

	if r != 0 {
		return err
	}

	return nil
}

func createClient(connectionStrings []string, iid string, p unsafe.Pointer) error {
	clzid, err := ole.IIDFromString(iid)
	if err != nil {
		return err
	}

	if len(connectionStrings) == 0 {
		return fmt.Errorf("empty connection string")
	}

	var conn []*uint16

	for _, c := range connectionStrings {
		s, err := windows.UTF16PtrFromString(c)

		if err != nil {
			return err
		}

		conn = append(conn, s)
	}

	r, _, err := fabricCreateClientProc.Call(
		uintptr(len(conn)),
		uintptr(unsafe.Pointer(&conn[0])),
		uintptr(unsafe.Pointer(clzid)),
		uintptr(p),
	)

	if r != 0 {
		return err
	}

	return nil
}

type FabricClient struct {
	createComObject func(iid string, outptr unsafe.Pointer) error
}

const (
	comIFabricClientSettingsIID = "{b0e7dee0-cf64-11e0-9572-0800200c9a66}" // Lowest ver Service Fabric 6.0
)

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

func NewLocalClient() (*FabricClient, error) {

	var com *ComFabricClientSettings
	err := createLocalClient(comIFabricClientSettingsIID, unsafe.Pointer(&com))
	if err != nil {
		return nil, err
	}

	return &FabricClient{func(iid string, outptr unsafe.Pointer) error {
		return queryObject(&com.IUnknown, iid, outptr)
	}}, nil
}

func NewInsecureClient(conn string) (*FabricClient, error) {
	var com *ComFabricClientSettings
	err := createClient([]string{conn}, comIFabricClientSettingsIID, unsafe.Pointer(&com))
	if err != nil {
		return nil, err
	}

	return &FabricClient{func(iid string, outptr unsafe.Pointer) error {
		return queryObject(&com.IUnknown, iid, outptr)
	}}, nil
}

func NewX509Client(conn string, cred X509Credentials) (*FabricClient, error) {
	var com *ComFabricClientSettings
	err := createClient([]string{conn}, comIFabricClientSettingsIID, unsafe.Pointer(&com))
	if err != nil {
		return nil, err
	}

	err = com.innerSetSecurityCredentials(&cred)

	if err != nil {
		return nil, err
	}

	return &FabricClient{func(iid string, outptr unsafe.Pointer) error {
		return queryObject(&com.IUnknown, iid, outptr)
	}}, nil
}
