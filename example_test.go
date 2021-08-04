package fabric_test

import (
	"context"
	"fmt"

	ole "github.com/go-ole/go-ole"
	"github.com/tg123/fabric"
)

func ExampleNewX509Client() {
	err := ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
	if err != nil {
		panic(err)
	}

	client, err := fabric.NewClient(fabric.FabricClientOpt{
		Address: []string{"test.southcentralus.cloudapp.azure.com:19000"},
		Credentials: &fabric.FabricSecurityCredentials{
			Kind: fabric.FabricSecurityCredentialKindX509,
			Value: fabric.FabricX509Credentials{
				FindType:              fabric.FabricX509FindTypeFindbythumbprint,
				FindValue:             "1111111111111111111111111111111111111111",
				StoreName:             "MY",
				StoreLocation:         fabric.FabricX509StoreLocationCurrentuser,
				RemoteCertThumbprints: []string{"1111111111111111111111111111111111111111"},
			},
		},
	})

	if err != nil {
		panic(err)
	}

	nodes, err := client.GetNodeList(context.TODO(), &fabric.FabricNodeQueryDescription{
		NodeNameFilter: "",
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(nodes)
}
