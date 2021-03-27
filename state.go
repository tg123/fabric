package fabric

import ole "github.com/go-ole/go-ole"

type ServiceContext struct {
	ServiceTypeName    string
	ServiceName        string
	InitializationData []byte
	PartitionId        ole.GUID
	InstanceId         int64
}
