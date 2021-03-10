package fabric

import (
	"fmt"
	"unsafe"

	ole "github.com/go-ole/go-ole"
	"golang.org/x/sys/windows"
)

var (
	fabricClientDll             = windows.NewLazyDLL("FabricClient.dll")
	fabricCreateLocalClientProc = fabricClientDll.NewProc("FabricCreateLocalClient")
	fabricCreateClientProc      = fabricClientDll.NewProc("FabricCreateClient")
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
	hub *fabricClientComHub
}

const (
	comIFabricClientSettingsIID = "{b0e7dee0-cf64-11e0-9572-0800200c9a66}" // Lowest ver Service Fabric 6.0
)

func clientFromComClientSetting(com *comFabricClientSettings) *FabricClient {
	hub := &fabricClientComHub{}
	hub.init(func(iid string, outptr unsafe.Pointer) error {
		return queryObject(&com.IUnknown, iid, outptr)
	})
	return &FabricClient{hub}
}

func NewLocalClient() (*FabricClient, error) {
	var com *comFabricClientSettings
	err := createLocalClient(comIFabricClientSettingsIID, unsafe.Pointer(&com))
	if err != nil {
		return nil, err
	}

	return clientFromComClientSetting(com), nil
}

func NewInsecureClient(conn string) (*FabricClient, error) {
	var com *comFabricClientSettings
	err := createClient([]string{conn}, comIFabricClientSettingsIID, unsafe.Pointer(&com))
	if err != nil {
		return nil, err
	}
	return clientFromComClientSetting(com), nil
}

func NewX509Client(conn string, cred X509Credentials) (*FabricClient, error) {
	var com *comFabricClientSettings
	err := createClient([]string{conn}, comIFabricClientSettingsIID, unsafe.Pointer(&com))
	if err != nil {
		return nil, err
	}

	err = com.innerSetSecurityCredentials(&cred)

	if err != nil {
		return nil, err
	}
	return clientFromComClientSetting(com), nil
}
