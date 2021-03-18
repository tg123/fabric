package fabric

import (
	"fmt"
	"sync"
	"syscall"
	"time"
	"unsafe"

	ole "github.com/go-ole/go-ole"
	"golang.org/x/sys/windows"
)

var (
	fabricClientDll              = windows.NewLazyDLL("FabricClient.dll")
	fabricCreateLocalClientProc  = fabricClientDll.NewProc("FabricCreateLocalClient")
	fabricCreateLocalClient2Proc = fabricClientDll.NewProc("FabricCreateLocalClient2")
	fabricCreateLocalClient3Proc = fabricClientDll.NewProc("FabricCreateLocalClient3")
	fabricCreateLocalClient4Proc = fabricClientDll.NewProc("FabricCreateLocalClient4")
	fabricCreateClientProc       = fabricClientDll.NewProc("FabricCreateClient")
	fabricCreateClient2Proc      = fabricClientDll.NewProc("FabricCreateClient2")
	fabricCreateClient3Proc      = fabricClientDll.NewProc("FabricCreateClient3")
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

func createClient(client *FabricClient, connectionStrings []string, iid string, p unsafe.Pointer) error {
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

	r, _, err := fabricCreateClient3Proc.Call(
		uintptr(len(conn)),
		uintptr(unsafe.Pointer(&conn[0])),
		uintptr(unsafe.Pointer(newFabricServiceNotificationEventHandler(client))),
		uintptr(unsafe.Pointer(newFabricConnectionEventHandler(client))),
		uintptr(unsafe.Pointer(clzid)),
		uintptr(p),
	)

	if r != 0 {
		return err
	}

	return nil
}

type goFabricServiceNotificationEventHandler struct {
	vtbl       *goFabricServiceNotificationEventHandlerVtbl
	unknownref *goIUnknown
	client     *FabricClient
}

type goFabricServiceNotificationEventHandlerVtbl struct {
	goIUnknownVtbl
	OnNotification uintptr
}

func (h *goFabricServiceNotificationEventHandler) onNotification(this *ole.IUnknown, notification *comFabricServiceNotification) uintptr {
	cb := h.client.OnNotification
	if cb == nil {
		return ole.S_OK
	}

	if notification == nil {
		return ole.E_POINTER
	}

	notification.AddRef()
	defer notification.Release()

	n, err := notification.GetNotification()
	if err != nil {
		return ole.E_FAIL
	}

	cb(*n)
	return ole.S_OK
}

func newFabricServiceNotificationEventHandler(client *FabricClient) *goFabricServiceNotificationEventHandler {
	handler := &goFabricServiceNotificationEventHandler{}
	handler.vtbl = &goFabricServiceNotificationEventHandlerVtbl{}
	handler.unknownref = attachIUnknown("{A04B7E9A-DAAB-45D4-8DA3-95EF3AB5DBAC}", &handler.vtbl.goIUnknownVtbl)
	handler.vtbl.OnNotification = syscall.NewCallback(handler.onNotification)
	handler.client = client
	client.deferclose = append(client.deferclose, func() {
		handler.unknownref.release((*ole.IUnknown)(unsafe.Pointer(handler)))
	})
	return handler
}

type goFabricConnectionEventHandler struct {
	vtbl       *goFabricConnectionEventHandlerVtbl
	unknownref *goIUnknown
	client     *FabricClient
}

type goFabricConnectionEventHandlerVtbl struct {
	goIUnknownVtbl
	OnConnected    uintptr
	OnDisconnected uintptr
}

func (h *goFabricConnectionEventHandler) onConnected(this *ole.IUnknown, result *comFabricGatewayInformationResult) uintptr {
	return h.onInfo(result, h.client.OnConnected)
}

func (h *goFabricConnectionEventHandler) onDisconnected(this *ole.IUnknown, result *comFabricGatewayInformationResult) uintptr {
	return h.onInfo(result, h.client.OnDisconnected)
}

func (h *goFabricConnectionEventHandler) onInfo(result *comFabricGatewayInformationResult, cb func(FabricGatewayInformation)) uintptr {
	if cb == nil {
		return ole.S_OK
	}

	if result == nil {
		return ole.E_POINTER
	}

	result.AddRef()
	defer result.Release()

	info, err := result.GetGatewayInformation()
	if err != nil {
		return ole.E_FAIL
	}

	cb(*info)
	return ole.S_OK
}

func newFabricConnectionEventHandler(client *FabricClient) *goFabricConnectionEventHandler {
	handler := &goFabricConnectionEventHandler{}
	handler.vtbl = &goFabricConnectionEventHandlerVtbl{}
	handler.unknownref = attachIUnknown("{2BD21F94-D962-4BB4-84B8-5A4B3E9D4D4D}", &handler.vtbl.goIUnknownVtbl)
	handler.vtbl.OnConnected = syscall.NewCallback(handler.onConnected)
	handler.vtbl.OnDisconnected = syscall.NewCallback(handler.onDisconnected)
	handler.client = client
	client.deferclose = append(client.deferclose, func() {
		handler.unknownref.release((*ole.IUnknown)(unsafe.Pointer(handler)))
	})
	return handler
}

type FabricClient struct {
	OnNotification func(notification FabricServiceNotification)
	OnConnected    func(info FabricGatewayInformation)
	OnDisconnected func(info FabricGatewayInformation)

	hub            *fabricClientComHub
	defaultTimeout time.Duration

	closed    bool
	closelock sync.Mutex

	deferclose []func()
}

func (v *FabricClient) GetTimeout() time.Duration {
	return v.defaultTimeout
}

func (v *FabricClient) SetDefaultTimeout(t time.Duration) {
	v.defaultTimeout = t
}

func (v *FabricClient) Close() error {
	// TODO calling func to a closing client is undefined
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

const (
	comIFabricClientSettingsIID = "{b0e7dee0-cf64-11e0-9572-0800200c9a66}" // Lowest ver Service Fabric 6.0
)

func comHubFromComClientSetting(com *comFabricClientSettings) *fabricClientComHub {
	hub := &fabricClientComHub{}
	hub.init(func(iid string, outptr unsafe.Pointer) error {
		return queryComObject(&com.IUnknown, iid, outptr)
	})
	releaseComObject(&com.IUnknown)
	return hub
}

// FabricClientOpt is the option to create a FabricClient
type FabricClientOpt struct {
	Address     []string
	Role        FabricClientRole
	Credentials *X509Credentials

	OnNotification func(notification FabricServiceNotification)
	OnConnected    func(info FabricGatewayInformation)
	OnDisconnected func(info FabricGatewayInformation)
}

func NewClient(opt FabricClientOpt) (*FabricClient, error) {
	c := &FabricClient{
		OnNotification: opt.OnNotification,
		OnConnected:    opt.OnConnected,
		OnDisconnected: opt.OnDisconnected,
		defaultTimeout: 5 * time.Minute, // opt.DefaultTimeout,
	}

	var com *comFabricClientSettings
	err := createClient(c, opt.Address, comIFabricClientSettingsIID, unsafe.Pointer(&com))
	if err != nil {
		return nil, err
	}

	if opt.Credentials != nil {
		err = com.innerSetSecurityCredentials(opt.Credentials)
		if err != nil {
			return nil, err
		}
	}

	hub := comHubFromComClientSetting(com)
	c.hub = hub

	return c, nil
}

func NewLocalClient() (*FabricClient, error) {
	var com *comFabricClientSettings
	err := createLocalClient(comIFabricClientSettingsIID, unsafe.Pointer(&com))
	if err != nil {
		return nil, err
	}

	// TODO support handlers
	c := &FabricClient{
		defaultTimeout: 5 * time.Minute,
	}
	hub := comHubFromComClientSetting(com)
	c.hub = hub

	return c, nil
}
