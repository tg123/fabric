package fabric

import (
	"fmt"
	"time"
	"unsafe"

	ole "github.com/go-ole/go-ole"
)

func createFabricSettingLocalClient(client *FabricClient, role FabricClientRole, p unsafe.Pointer) error {
	err := callCreateLocalClient4(
		unsafe.Pointer(newComFabricServiceNotificationEventHandler(client)),
		unsafe.Pointer(newComFabricClientConnectionEventHandler(client)),
		role,
		unsafe.Pointer(iidIFabricClientSettingsIID),
		p,
	)

	if err != nil {
		return err
	}

	return nil
}

func createFabricSettingClient(client *FabricClient, connectionStrings []string, p unsafe.Pointer) error {
	if len(connectionStrings) == 0 {
		return fmt.Errorf("empty connection string")
	}

	var conn []*uint16

	for _, c := range connectionStrings {
		conn = append(conn, utf16PtrFromString(c))
	}

	err := callCreateClient3(
		len(conn),
		unsafe.Pointer(&conn[0]),
		unsafe.Pointer(newComFabricServiceNotificationEventHandler(client)),
		unsafe.Pointer(newComFabricClientConnectionEventHandler(client)),
		unsafe.Pointer(iidIFabricClientSettingsIID),
		p,
	)

	if err != nil {
		return err
	}

	return nil
}

func (h *comFabricServiceNotificationEventHandlerGoProxy) init() {
	h.client.deferclose = append(h.client.deferclose, func() {
		h.unknownref.release((*ole.IUnknown)(unsafe.Pointer(h)))
	})
}

func (h *comFabricServiceNotificationEventHandlerGoProxy) OnNotification(this *ole.IUnknown, notification *comFabricServiceNotification) uintptr {
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

func (h *comFabricClientConnectionEventHandlerGoProxy) init() {
	h.client.deferclose = append(h.client.deferclose, func() {
		h.unknownref.release((*ole.IUnknown)(unsafe.Pointer(h)))
	})
}

func (h *comFabricClientConnectionEventHandlerGoProxy) OnConnected(this *ole.IUnknown, result *comFabricGatewayInformationResult) uintptr {
	return h.OnInfo(result, h.client.OnConnected)
}

func (h *comFabricClientConnectionEventHandlerGoProxy) OnDisconnected(this *ole.IUnknown, result *comFabricGatewayInformationResult) uintptr {
	return h.OnInfo(result, h.client.OnDisconnected)
}

func (h *comFabricClientConnectionEventHandlerGoProxy) OnInfo(result *comFabricGatewayInformationResult, cb func(FabricGatewayInformation)) uintptr {
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

type FabricClient struct {
	coclz
	hub *fabricClientComHub

	OnNotification func(notification FabricServiceNotification)
	OnConnected    func(info FabricGatewayInformation)
	OnDisconnected func(info FabricGatewayInformation)
}

var (
	iidIFabricClientSettingsIID = ole.NewGUID("{B0E7DEE0-CF64-11E0-9572-0800200C9A66}") // Lowest ver Service Fabric 6.0
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
	Address        []string
	Credentials    *FabricSecurityCredentials
	OnNotification func(notification FabricServiceNotification)
	OnConnected    func(info FabricGatewayInformation)
	OnDisconnected func(info FabricGatewayInformation)
}

type FabricLocalClientOpt struct {
	Role           FabricClientRole
	OnNotification func(notification FabricServiceNotification)
	OnConnected    func(info FabricGatewayInformation)
	OnDisconnected func(info FabricGatewayInformation)
}

func NewClient(opt FabricClientOpt) (*FabricClient, error) {

	c := &FabricClient{
		OnNotification: opt.OnNotification,
		OnConnected:    opt.OnConnected,
		OnDisconnected: opt.OnDisconnected,
	}

	var com *comFabricClientSettings
	err := createFabricSettingClient(c, opt.Address, unsafe.Pointer(&com))
	if err != nil {
		return nil, err
	}

	if opt.Credentials != nil {
		err = com.SetSecurityCredentials(opt.Credentials)
		if err != nil {
			return nil, err
		}
	}

	hub := comHubFromComClientSetting(com)
	c.hub = hub
	c.coclz.hub = hub
	c.coclz.defaultTimeout = 5 * time.Minute

	return c, nil
}

func NewInsecureClient(conn string) (*FabricClient, error) {
	return NewClient(FabricClientOpt{
		Address: []string{conn},
	})
}

func NewX509Client(conn string, cred *FabricSecurityCredentials) (*FabricClient, error) {
	return NewClient(FabricClientOpt{
		Address: []string{conn},
		Credentials: &FabricSecurityCredentials{
			Kind:  FabricSecurityCredentialKindX509,
			Value: cred,
		},
	})
}

func NewLocalClient() (*FabricClient, error) {
	return NewLocalClientOpt(FabricLocalClientOpt{
		Role: FabricClientRoleAdmin,
	})
}

func NewLocalClientOpt(opt FabricLocalClientOpt) (*FabricClient, error) {
	c := &FabricClient{
		OnNotification: opt.OnNotification,
		OnConnected:    opt.OnConnected,
		OnDisconnected: opt.OnDisconnected,
	}

	var com *comFabricClientSettings
	err := createFabricSettingLocalClient(c, opt.Role, unsafe.Pointer(&com))
	if err != nil {
		return nil, err
	}

	hub := comHubFromComClientSetting(com)
	c.hub = hub
	c.coclz.hub = hub
	c.coclz.defaultTimeout = 5 * time.Minute

	return c, nil
}
