# Service Fabric Golang SDK (COM+)

[![](https://pkg.go.dev/badge/github.com/tg123/fabric?status.svg)](https://pkg.go.dev/github.com/tg123/fabric)

## status: UNOFFICIAL and WIP

This is SDK is generated from IDL file in service fabric repo <https://github.com/microsoft/service-fabric/tree/master/src/prod/src/idl/public>.
The package calls Service Fabric COM+ API directly which is the same to what dotnet SDK does.

Working in process and will provide friendly programming experience to Golang users.

Current features:

 * BeginXX and EndXX are combined into goroutine style
 * common data types, for example `FILETIME`, are mapped into go types

## Usage

make sure Service Fabric installed and be visible in your PATH

### Example

```
package main

import (
	"fmt"
	"time"

	ole "github.com/go-ole/go-ole"
	"github.com/tg123/fabric"
)

func main() {
	err := ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
	if err != nil {
		panic(err)
	}

	client, err := fabric.NewX509Client("test.southcentralus.cloudapp.azure.com:19000", fabric.X509Credentials{
		FindType:              fabric.FabricX509FindTypeFindbythumbprint,
		FindValue:             "1111111111111111111111111111111111111111",
		StoreName:             "MY",
		StoreLocation:         fabric.FabricX509StoreLocationCurrentuser,
		RemoteCertThumbprints: []string{"1111111111111111111111111111111111111111"},
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
```
