// Code generated by "go run github.com/tg123/fabric/mkidl"; DO NOT EDIT.
package fabric

import (
	"unsafe"
)

type fabricCodePackageActivationContextComHub struct {
	FabricCodePackageActivationContext  *comFabricCodePackageActivationContext
	FabricCodePackageActivationContext2 *comFabricCodePackageActivationContext2
	FabricCodePackageActivationContext3 *comFabricCodePackageActivationContext3
	FabricCodePackageActivationContext4 *comFabricCodePackageActivationContext4
	FabricCodePackageActivationContext5 *comFabricCodePackageActivationContext5
	FabricCodePackageActivationContext6 *comFabricCodePackageActivationContext6
}

func (v *fabricCodePackageActivationContextComHub) init(createComObject comCreator) {
	createComObject("{68A971E2-F15F-4D95-A79C-8A257909659E}", unsafe.Pointer(&v.FabricCodePackageActivationContext))
	createComObject("{6C83D5C1-1954-4B80-9175-0D0E7C8715C9}", unsafe.Pointer(&v.FabricCodePackageActivationContext2))
	createComObject("{6EFEE900-F491-4B03-BC5B-3A70DE103593}", unsafe.Pointer(&v.FabricCodePackageActivationContext3))
	createComObject("{99EFEBB6-A7B4-4D45-B45E-F191A66EEF03}", unsafe.Pointer(&v.FabricCodePackageActivationContext4))
	createComObject("{FE45387E-8711-4949-AC36-31DC95035513}", unsafe.Pointer(&v.FabricCodePackageActivationContext5))
	createComObject("{FA5FDA9B-472C-45A0-9B60-A374691227A4}", unsafe.Pointer(&v.FabricCodePackageActivationContext6))
}
func (v *fabricCodePackageActivationContextComHub) Close() error {
	if v.FabricCodePackageActivationContext != nil {
		releaseComObject(&v.FabricCodePackageActivationContext.IUnknown)
	}
	if v.FabricCodePackageActivationContext2 != nil {
		releaseComObject(&v.FabricCodePackageActivationContext2.IUnknown)
	}
	if v.FabricCodePackageActivationContext3 != nil {
		releaseComObject(&v.FabricCodePackageActivationContext3.IUnknown)
	}
	if v.FabricCodePackageActivationContext4 != nil {
		releaseComObject(&v.FabricCodePackageActivationContext4.IUnknown)
	}
	if v.FabricCodePackageActivationContext5 != nil {
		releaseComObject(&v.FabricCodePackageActivationContext5.IUnknown)
	}
	if v.FabricCodePackageActivationContext6 != nil {
		releaseComObject(&v.FabricCodePackageActivationContext6.IUnknown)
	}
	return nil
}
func (v *FabricCodePackageActivationContext) GetContextId() (rt string, err error) {
	if v.hub.FabricCodePackageActivationContext == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext.GetContextId()
}
func (v *FabricCodePackageActivationContext) GetCodePackageName() (rt string, err error) {
	if v.hub.FabricCodePackageActivationContext == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext.GetCodePackageName()
}
func (v *FabricCodePackageActivationContext) GetCodePackageVersion() (rt string, err error) {
	if v.hub.FabricCodePackageActivationContext == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext.GetCodePackageVersion()
}
func (v *FabricCodePackageActivationContext) GetWorkDirectory() (rt string, err error) {
	if v.hub.FabricCodePackageActivationContext == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext.GetWorkDirectory()
}
func (v *FabricCodePackageActivationContext) GetLogDirectory() (rt string, err error) {
	if v.hub.FabricCodePackageActivationContext == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext.GetLogDirectory()
}
func (v *FabricCodePackageActivationContext) GetTempDirectory() (rt string, err error) {
	if v.hub.FabricCodePackageActivationContext == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext.GetTempDirectory()
}
func (v *FabricCodePackageActivationContext) GetServiceTypes() (rt []FabricServiceTypeDescription, err error) {
	if v.hub.FabricCodePackageActivationContext == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext.GetServiceTypes()
}
func (v *FabricCodePackageActivationContext) GetServiceGroupTypes() (rt []FabricServiceGroupTypeDescription, err error) {
	if v.hub.FabricCodePackageActivationContext == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext.GetServiceGroupTypes()
}
func (v *FabricCodePackageActivationContext) GetApplicationPrincipals() (rt *FabricApplicationPrincipalsDescription, err error) {
	if v.hub.FabricCodePackageActivationContext == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext.GetApplicationPrincipals()
}
func (v *FabricCodePackageActivationContext) GetServiceEndpointResources() (rt []FabricEndpointResourceDescription, err error) {
	if v.hub.FabricCodePackageActivationContext == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext.GetServiceEndpointResources()
}
func (v *FabricCodePackageActivationContext) GetServiceEndpointResource(
	serviceEndpointResourceName string,
) (bufferedValue *FabricEndpointResourceDescription, err error) {
	if v.hub.FabricCodePackageActivationContext == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext.GetServiceEndpointResource(serviceEndpointResourceName)
}
func (v *FabricCodePackageActivationContext) GetCodePackageNames() (names *comFabricStringListResult, err error) {
	if v.hub.FabricCodePackageActivationContext == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext.GetCodePackageNames()
}
func (v *FabricCodePackageActivationContext) GetConfigurationPackageNames() (names *comFabricStringListResult, err error) {
	if v.hub.FabricCodePackageActivationContext == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext.GetConfigurationPackageNames()
}
func (v *FabricCodePackageActivationContext) GetDataPackageNames() (names *comFabricStringListResult, err error) {
	if v.hub.FabricCodePackageActivationContext == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext.GetDataPackageNames()
}
func (v *FabricCodePackageActivationContext) GetCodePackage(
	codePackageName string,
) (codePackage *comFabricCodePackage, err error) {
	if v.hub.FabricCodePackageActivationContext == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext.GetCodePackage(codePackageName)
}
func (v *FabricCodePackageActivationContext) GetConfigurationPackage(
	configPackageName string,
) (configPackage *comFabricConfigurationPackage, err error) {
	if v.hub.FabricCodePackageActivationContext == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext.GetConfigurationPackage(configPackageName)
}
func (v *FabricCodePackageActivationContext) GetDataPackage(
	dataPackageName string,
) (dataPackage *comFabricDataPackage, err error) {
	if v.hub.FabricCodePackageActivationContext == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext.GetDataPackage(dataPackageName)
}
func (v *FabricCodePackageActivationContext) registerCodePackageChangeHandler(
	callback *comFabricCodePackageChangeHandler,
) (callbackHandle int64, err error) {
	if v.hub.FabricCodePackageActivationContext == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext.RegisterCodePackageChangeHandler(callback)
}
func (v *FabricCodePackageActivationContext) unregisterCodePackageChangeHandler(
	callbackHandle int64,
) (err error) {
	if v.hub.FabricCodePackageActivationContext == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext.UnregisterCodePackageChangeHandler(callbackHandle)
}
func (v *FabricCodePackageActivationContext) registerConfigurationPackageChangeHandler(
	callback *comFabricConfigurationPackageChangeHandler,
) (callbackHandle int64, err error) {
	if v.hub.FabricCodePackageActivationContext == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext.RegisterConfigurationPackageChangeHandler(callback)
}
func (v *FabricCodePackageActivationContext) unregisterConfigurationPackageChangeHandler(
	callbackHandle int64,
) (err error) {
	if v.hub.FabricCodePackageActivationContext == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext.UnregisterConfigurationPackageChangeHandler(callbackHandle)
}
func (v *FabricCodePackageActivationContext) registerDataPackageChangeHandler(
	callback *comFabricDataPackageChangeHandler,
) (callbackHandle int64, err error) {
	if v.hub.FabricCodePackageActivationContext == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext.RegisterDataPackageChangeHandler(callback)
}
func (v *FabricCodePackageActivationContext) unregisterDataPackageChangeHandler(
	callbackHandle int64,
) (err error) {
	if v.hub.FabricCodePackageActivationContext == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext.UnregisterDataPackageChangeHandler(callbackHandle)
}
func (v *FabricCodePackageActivationContext) GetApplicationName() (rt string, err error) {
	if v.hub.FabricCodePackageActivationContext2 == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext2.GetApplicationName()
}
func (v *FabricCodePackageActivationContext) GetApplicationTypeName() (rt string, err error) {
	if v.hub.FabricCodePackageActivationContext2 == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext2.GetApplicationTypeName()
}
func (v *FabricCodePackageActivationContext) GetServiceManifestName() (serviceManifestName *comFabricStringResult, err error) {
	if v.hub.FabricCodePackageActivationContext2 == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext2.GetServiceManifestName()
}
func (v *FabricCodePackageActivationContext) GetServiceManifestVersion() (serviceManifestVersion *comFabricStringResult, err error) {
	if v.hub.FabricCodePackageActivationContext2 == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext2.GetServiceManifestVersion()
}
func (v *FabricCodePackageActivationContext) ReportApplicationHealth(
	healthInfo *FabricHealthInformation,
) (err error) {
	if v.hub.FabricCodePackageActivationContext3 == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext3.ReportApplicationHealth(healthInfo)
}
func (v *FabricCodePackageActivationContext) ReportDeployedApplicationHealth(
	healthInfo *FabricHealthInformation,
) (err error) {
	if v.hub.FabricCodePackageActivationContext3 == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext3.ReportDeployedApplicationHealth(healthInfo)
}
func (v *FabricCodePackageActivationContext) ReportDeployedServicePackageHealth(
	healthInfo *FabricHealthInformation,
) (err error) {
	if v.hub.FabricCodePackageActivationContext3 == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext3.ReportDeployedServicePackageHealth(healthInfo)
}
func (v *FabricCodePackageActivationContext) ReportApplicationHealth2(
	healthInfo *FabricHealthInformation,
	sendOptions *FabricHealthReportSendOptions,
) (err error) {
	if v.hub.FabricCodePackageActivationContext4 == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext4.ReportApplicationHealth2(healthInfo, sendOptions)
}
func (v *FabricCodePackageActivationContext) ReportDeployedApplicationHealth2(
	healthInfo *FabricHealthInformation,
	sendOptions *FabricHealthReportSendOptions,
) (err error) {
	if v.hub.FabricCodePackageActivationContext4 == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext4.ReportDeployedApplicationHealth2(healthInfo, sendOptions)
}
func (v *FabricCodePackageActivationContext) ReportDeployedServicePackageHealth2(
	healthInfo *FabricHealthInformation,
	sendOptions *FabricHealthReportSendOptions,
) (err error) {
	if v.hub.FabricCodePackageActivationContext4 == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext4.ReportDeployedServicePackageHealth2(healthInfo, sendOptions)
}
func (v *FabricCodePackageActivationContext) GetServiceListenAddress() (rt string, err error) {
	if v.hub.FabricCodePackageActivationContext5 == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext5.GetServiceListenAddress()
}
func (v *FabricCodePackageActivationContext) GetServicePublishAddress() (rt string, err error) {
	if v.hub.FabricCodePackageActivationContext5 == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext5.GetServicePublishAddress()
}
func (v *FabricCodePackageActivationContext) GetDirectory(
	logicalDirectoryName string,
) (directoryPath *comFabricStringResult, err error) {
	if v.hub.FabricCodePackageActivationContext6 == nil {
		err = errComNotImpl
		return
	}
	return v.hub.FabricCodePackageActivationContext6.GetDirectory(logicalDirectoryName)
}