// +build generate

package fabric

//go:generate go run github.com/tg123/fabric/mkidl idls/FabricCommon.idl idls/FabricTypes.idl idls/FabricClient.idl idls/FabricRuntime.idl
//go:generate go fmt