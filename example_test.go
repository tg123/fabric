package fabric

import (
	"context"
	"fmt"

	ole "github.com/go-ole/go-ole"
)

func ExampleNewX509Client() {
	err := ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
	if err != nil {
		panic(err)
	}

	client, err := NewX509Client("test.southcentralus.cloudapp.azure.com:19000", X509Credentials{
		FindType:              FabricX509FindTypeFindbythumbprint,
		FindValue:             "1111111111111111111111111111111111111111",
		StoreName:             "MY",
		StoreLocation:         FabricX509StoreLocationCurrentuser,
		RemoteCertThumbprints: []string{"1111111111111111111111111111111111111111"},
	})

	if err != nil {
		panic(err)
	}

	q, err := client.CreateFabricQueryClient8()
	if err != nil {
		panic(err)
	}

	nodes, err := q.GetNodeList(context.TODO(), &FabricNodeQueryDescription{
		NodeNameFilter: "",
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(nodes)
}