package fabric

import "github.com/google/uuid"

type ServiceContext struct {
	ServiceTypeName    string
	ServiceName        string
	InitializationData []byte
	PartitionId        uuid.UUID
	InstanceId         int64
}
