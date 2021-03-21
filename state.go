package fabric

import "golang.org/x/sys/windows"

type ServiceContext struct {
	ServiceTypeName    string
	ServiceName        string
	InitializationData []byte
	PartitionId        windows.GUID
	InstanceId         int64
}
