// Code generated by "stringer -type FabricErrorCode"; DO NOT EDIT.

package fabric

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[FabricErrorFirstReservedHresult-2147949500]
	_ = x[FabricErrorLastReservedHresult-2147949899]
	_ = x[FabricErrorCommunicationError-2147949500]
	_ = x[FabricErrorInvalidAddress-2147949501]
	_ = x[FabricErrorInvalidNameUri-2147949502]
	_ = x[FabricErrorInvalidPartitionKey-2147949503]
	_ = x[FabricErrorNameAlreadyExists-2147949504]
	_ = x[FabricErrorNameDoesNotExist-2147949505]
	_ = x[FabricErrorNameNotEmpty-2147949506]
	_ = x[FabricErrorNodeNotFound-2147949507]
	_ = x[FabricErrorNodeIsUp-2147949508]
	_ = x[FabricErrorNoWriteQuorum-2147949509]
	_ = x[FabricErrorNotPrimary-2147949510]
	_ = x[FabricErrorNotReady-2147949511]
	_ = x[FabricErrorOperationNotComplete-2147949512]
	_ = x[FabricErrorPropertyDoesNotExist-2147949513]
	_ = x[FabricErrorReconfigurationPending-2147949514]
	_ = x[FabricErrorReplicationQueueFull-2147949515]
	_ = x[FabricErrorServiceAlreadyExists-2147949516]
	_ = x[FabricErrorServiceDoesNotExist-2147949517]
	_ = x[FabricErrorServiceOffline-2147949518]
	_ = x[FabricErrorServiceMetadataMismatch-2147949519]
	_ = x[FabricErrorServiceAffinityChainNotSupported-2147949520]
	_ = x[FabricErrorServiceTypeAlreadyRegistered-2147949521]
	_ = x[FabricErrorServiceTypeNotRegistered-2147949522]
	_ = x[FabricErrorValueTooLarge-2147949523]
	_ = x[FabricErrorValueEmpty-2147949524]
	_ = x[FabricErrorPropertyCheckFailed-2147949525]
	_ = x[FabricErrorWriteConflict-2147949526]
	_ = x[FabricErrorEnumerationCompleted-2147949527]
	_ = x[FabricErrorApplicationTypeProvisionInProgress-2147949528]
	_ = x[FabricErrorApplicationTypeAlreadyExists-2147949529]
	_ = x[FabricErrorApplicationTypeNotFound-2147949530]
	_ = x[FabricErrorApplicationTypeInUse-2147949531]
	_ = x[FabricErrorApplicationAlreadyExists-2147949532]
	_ = x[FabricErrorApplicationNotFound-2147949533]
	_ = x[FabricErrorApplicationUpgradeInProgress-2147949534]
	_ = x[FabricErrorApplicationUpgradeValidationError-2147949535]
	_ = x[FabricErrorServiceTypeNotFound-2147949536]
	_ = x[FabricErrorServiceTypeMismatch-2147949537]
	_ = x[FabricErrorServiceTypeTemplateNotFound-2147949538]
	_ = x[FabricErrorConfigurationSectionNotFound-2147949539]
	_ = x[FabricErrorConfigurationParameterNotFound-2147949540]
	_ = x[FabricErrorInvalidConfiguration-2147949541]
	_ = x[FabricErrorImagebuilderValidationError-2147949542]
	_ = x[FabricErrorPartitionNotFound-2147949543]
	_ = x[FabricErrorReplicaDoesNotExist-2147949544]
	_ = x[FabricErrorServiceGroupAlreadyExists-2147949545]
	_ = x[FabricErrorServiceGroupDoesNotExist-2147949546]
	_ = x[FabricErrorProcessDeactivated-2147949547]
	_ = x[FabricErrorProcessAborted-2147949548]
	_ = x[FabricErrorUpgradeFailed-2147949549]
	_ = x[FabricErrorInvalidCredentialType-2147949550]
	_ = x[FabricErrorInvalidX509FindType-2147949551]
	_ = x[FabricErrorInvalidX509StoreLocation-2147949552]
	_ = x[FabricErrorInvalidX509StoreName-2147949553]
	_ = x[FabricErrorInvalidX509Thumbprint-2147949554]
	_ = x[FabricErrorInvalidProtectionLevel-2147949555]
	_ = x[FabricErrorInvalidX509Store-2147949556]
	_ = x[FabricErrorInvalidSubjectName-2147949557]
	_ = x[FabricErrorInvalidAllowedCommonNameList-2147949558]
	_ = x[FabricErrorInvalidCredentials-2147949559]
	_ = x[FabricErrorDecryptionFailed-2147949560]
	_ = x[FabricErrorConfigurationPackageNotFound-2147949561]
	_ = x[FabricErrorDataPackageNotFound-2147949562]
	_ = x[FabricErrorCodePackageNotFound-2147949563]
	_ = x[FabricErrorServiceEndpointResourceNotFound-2147949564]
	_ = x[FabricErrorInvalidOperation-2147949565]
	_ = x[FabricErrorObjectClosed-2147949566]
	_ = x[FabricErrorTimeout-2147949567]
	_ = x[FabricErrorFileNotFound-2147949568]
	_ = x[FabricErrorDirectoryNotFound-2147949569]
	_ = x[FabricErrorInvalidDirectory-2147949570]
	_ = x[FabricErrorPathTooLong-2147949571]
	_ = x[FabricErrorImagestoreIoerror-2147949572]
	_ = x[FabricErrorCorruptedImageStoreObjectFound-2147949573]
	_ = x[FabricErrorApplicationNotUpgrading-2147949574]
	_ = x[FabricErrorApplicationAlreadyInTargetVersion-2147949575]
	_ = x[FabricErrorImagebuilderUnexpectedError-2147949576]
	_ = x[FabricErrorFabricVersionNotFound-2147949577]
	_ = x[FabricErrorFabricVersionInUse-2147949578]
	_ = x[FabricErrorFabricVersionAlreadyExists-2147949579]
	_ = x[FabricErrorFabricAlreadyInTargetVersion-2147949580]
	_ = x[FabricErrorFabricNotUpgrading-2147949581]
	_ = x[FabricErrorFabricUpgradeInProgress-2147949582]
	_ = x[FabricErrorFabricUpgradeValidationError-2147949583]
	_ = x[FabricErrorHealthMaxReportsReached-2147949584]
	_ = x[FabricErrorHealthStaleReport-2147949585]
	_ = x[FabricErrorKeyTooLarge-2147949586]
	_ = x[FabricErrorKeyNotFound-2147949587]
	_ = x[FabricErrorSequenceNumberCheckFailed-2147949588]
	_ = x[FabricErrorEncryptionFailed-2147949589]
	_ = x[FabricErrorInvalidAtomicGroup-2147949590]
	_ = x[FabricErrorHealthEntityNotFound-2147949591]
	_ = x[FabricErrorServiceManifestNotFound-2147949592]
	_ = x[FabricErrorReliableSessionTransportStartupFailure-2147949593]
	_ = x[FabricErrorReliableSessionAlreadyExists-2147949594]
	_ = x[FabricErrorReliableSessionCannotConnect-2147949595]
	_ = x[FabricErrorReliableSessionManagerExists-2147949596]
	_ = x[FabricErrorReliableSessionRejected-2147949597]
	_ = x[FabricErrorReliableSessionManagerAlreadyListening-2147949598]
	_ = x[FabricErrorReliableSessionManagerNotFound-2147949599]
	_ = x[FabricErrorReliableSessionManagerNotListening-2147949600]
	_ = x[FabricErrorInvalidServiceType-2147949601]
	_ = x[FabricErrorImagebuilderTimeout-2147949602]
	_ = x[FabricErrorImagebuilderAccessDenied-2147949603]
	_ = x[FabricErrorImagebuilderInvalidMsiFile-2147949604]
	_ = x[FabricErrorServiceTooBusy-2147949605]
	_ = x[FabricErrorTransactionNotActive-2147949606]
	_ = x[FabricErrorRepairTaskAlreadyExists-2147949607]
	_ = x[FabricErrorRepairTaskNotFound-2147949608]
	_ = x[FabricErrorReliableSessionNotFound-2147949609]
	_ = x[FabricErrorReliableSessionQueueEmpty-2147949610]
	_ = x[FabricErrorReliableSessionQuotaExceeded-2147949611]
	_ = x[FabricErrorReliableSessionServiceFaulted-2147949612]
	_ = x[FabricErrorReliableSessionInvalidTargetPartition-2147949613]
	_ = x[FabricErrorTransactionTooLarge-2147949614]
	_ = x[FabricErrorReplicationOperationTooLarge-2147949615]
	_ = x[FabricErrorInstanceIdMismatch-2147949616]
	_ = x[FabricErrorUpgradeDomainAlreadyCompleted-2147949617]
	_ = x[FabricErrorNodeHasNotStoppedYet-2147949618]
	_ = x[FabricErrorInsufficientClusterCapacity-2147949619]
	_ = x[FabricErrorInvalidPackageSharingPolicy-2147949620]
	_ = x[FabricErrorPredeploymentNotAllowed-2147949621]
	_ = x[FabricErrorInvalidBackupSetting-2147949622]
	_ = x[FabricErrorMissingFullBackup-2147949623]
	_ = x[FabricErrorBackupInProgress-2147949624]
	_ = x[FabricErrorDuplicateServiceNotificationFilterName-2147949625]
	_ = x[FabricErrorInvalidReplicaOperation-2147949626]
	_ = x[FabricErrorInvalidReplicaState-2147949627]
	_ = x[FabricErrorLoadbalancerNotReady-2147949628]
	_ = x[FabricErrorInvalidPartitionOperation-2147949629]
	_ = x[FabricErrorPrimaryAlreadyExists-2147949630]
	_ = x[FabricErrorSecondaryAlreadyExists-2147949631]
	_ = x[FabricErrorBackupDirectoryNotEmpty-2147949632]
	_ = x[FabricErrorForceNotSupportedForReplicaOperation-2147949633]
	_ = x[FabricErrorAcquireFileLockFailed-2147949634]
	_ = x[FabricErrorConnectionDenied-2147949635]
	_ = x[FabricErrorServerAuthenticationFailed-2147949636]
	_ = x[FabricErrorConstraintKeyUndefined-2147949637]
	_ = x[FabricErrorMultithreadedTransactionsNotAllowed-2147949638]
	_ = x[FabricErrorInvalidX509NameList-2147949639]
	_ = x[FabricErrorVerboseFmPlacementHealthReportingRequired-2147949640]
	_ = x[FabricErrorGatewayNotReachable-2147949641]
	_ = x[FabricErrorUserRoleClientCertificateNotConfigured-2147949642]
	_ = x[FabricErrorTransactionAborted-2147949643]
	_ = x[FabricErrorCannotConnect-2147949644]
	_ = x[FabricErrorMessageTooLarge-2147949645]
	_ = x[FabricErrorConstraintNotSatisfied-2147949646]
	_ = x[FabricErrorEndpointNotFound-2147949647]
	_ = x[FabricErrorApplicationUpdateInProgress-2147949648]
	_ = x[FabricErrorDeleteBackupFileFailed-2147949649]
	_ = x[FabricErrorConnectionClosedByRemoteEnd-2147949650]
	_ = x[FabricErrorInvalidTestCommandState-2147949651]
	_ = x[FabricErrorTestCommandOperationIdAlreadyExists-2147949652]
	_ = x[FabricErrorCmOperationFailed-2147949653]
	_ = x[FabricErrorImagebuilderReservedDirectoryError-2147949654]
	_ = x[FabricErrorCertificateNotFound-2147949655]
	_ = x[FabricErrorChaosAlreadyRunning-2147949656]
	_ = x[FabricErrorFabricDataRootNotFound-2147949657]
	_ = x[FabricErrorInvalidRestoreData-2147949658]
	_ = x[FabricErrorDuplicateBackups-2147949659]
	_ = x[FabricErrorInvalidBackupChain-2147949660]
	_ = x[FabricErrorStopInProgress-2147949661]
	_ = x[FabricErrorAlreadyStopped-2147949662]
	_ = x[FabricErrorNodeIsDown-2147949663]
	_ = x[FabricErrorNodeTransitionInProgress-2147949664]
	_ = x[FabricErrorInvalidBackup-2147949665]
	_ = x[FabricErrorInvalidInstanceId-2147949666]
	_ = x[FabricErrorInvalidDuration-2147949667]
	_ = x[FabricErrorRestoreSafeCheckFailed-2147949668]
	_ = x[FabricErrorConfigUpgradeFailed-2147949669]
	_ = x[FabricErrorUploadSessionRangeNotSatisfiable-2147949670]
	_ = x[FabricErrorUploadSessionIdConflict-2147949671]
	_ = x[FabricErrorInvalidPartitionSelector-2147949672]
	_ = x[FabricErrorInvalidReplicaSelector-2147949673]
	_ = x[FabricErrorDnsServiceNotFound-2147949674]
	_ = x[FabricErrorInvalidDnsName-2147949675]
	_ = x[FabricErrorDnsNameInUse-2147949676]
	_ = x[FabricErrorComposeDeploymentAlreadyExists-2147949677]
	_ = x[FabricErrorComposeDeploymentNotFound-2147949678]
	_ = x[FabricErrorInvalidForStatefulServices-2147949679]
	_ = x[FabricErrorInvalidForStatelessServices-2147949680]
	_ = x[FabricErrorOnlyValidForStatefulPersistentServices-2147949681]
	_ = x[FabricErrorInvalidUploadSessionId-2147949682]
	_ = x[FabricErrorBackupNotEnabled-2147949683]
	_ = x[FabricErrorBackupIsEnabled-2147949684]
	_ = x[FabricErrorBackupPolicyDoesNotExist-2147949685]
	_ = x[FabricErrorBackupPolicyAlreadyExists-2147949686]
	_ = x[FabricErrorRestoreInProgress-2147949687]
	_ = x[FabricErrorRestoreSourceTargetPartitionMismatch-2147949688]
	_ = x[FabricErrorFaultAnalysisServiceNotEnabled-2147949689]
	_ = x[FabricErrorContainerNotFound-2147949690]
	_ = x[FabricErrorObjectDisposed-2147949691]
	_ = x[FabricErrorNotReadable-2147949692]
	_ = x[FabricErrorBackupcopierUnexpectedError-2147949693]
	_ = x[FabricErrorBackupcopierTimeout-2147949694]
	_ = x[FabricErrorBackupcopierAccessDenied-2147949695]
	_ = x[FabricErrorInvalidServiceScalingPolicy-2147949696]
	_ = x[FabricErrorSingleInstanceApplicationAlreadyExists-2147949697]
	_ = x[FabricErrorSingleInstanceApplicationNotFound-2147949698]
	_ = x[FabricErrorVolumeAlreadyExists-2147949699]
	_ = x[FabricErrorVolumeNotFound-2147949700]
	_ = x[FabricErrorDatabaseMigrationInProgress-2147949701]
	_ = x[FabricErrorCentralSecretServiceGeneric-2147949702]
	_ = x[FabricErrorSecretInvalid-2147949703]
	_ = x[FabricErrorSecretVersionAlreadyExists-2147949704]
	_ = x[FabricErrorSingleInstanceApplicationUpgradeInProgress-2147949705]
	_ = x[FabricErrorOperationNotSupported-2147949706]
	_ = x[FabricErrorComposeDeploymentNotUpgrading-2147949707]
	_ = x[FabricErrorSecretTypeCannotBeChanged-2147949708]
	_ = x[FabricErrorNetworkNotFound-2147949709]
	_ = x[FabricErrorNetworkInUse-2147949710]
	_ = x[FabricErrorEndpointNotReferenced-2147949711]
	_ = x[FabricErrorLastUsedHresult-2147949711]
}

const (
	_FabricErrorCode_name_0 = "FabricErrorFirstReservedHresultFabricErrorInvalidAddressFabricErrorInvalidNameUriFabricErrorInvalidPartitionKeyFabricErrorNameAlreadyExistsFabricErrorNameDoesNotExistFabricErrorNameNotEmptyFabricErrorNodeNotFoundFabricErrorNodeIsUpFabricErrorNoWriteQuorumFabricErrorNotPrimaryFabricErrorNotReadyFabricErrorOperationNotCompleteFabricErrorPropertyDoesNotExistFabricErrorReconfigurationPendingFabricErrorReplicationQueueFullFabricErrorServiceAlreadyExistsFabricErrorServiceDoesNotExistFabricErrorServiceOfflineFabricErrorServiceMetadataMismatchFabricErrorServiceAffinityChainNotSupportedFabricErrorServiceTypeAlreadyRegisteredFabricErrorServiceTypeNotRegisteredFabricErrorValueTooLargeFabricErrorValueEmptyFabricErrorPropertyCheckFailedFabricErrorWriteConflictFabricErrorEnumerationCompletedFabricErrorApplicationTypeProvisionInProgressFabricErrorApplicationTypeAlreadyExistsFabricErrorApplicationTypeNotFoundFabricErrorApplicationTypeInUseFabricErrorApplicationAlreadyExistsFabricErrorApplicationNotFoundFabricErrorApplicationUpgradeInProgressFabricErrorApplicationUpgradeValidationErrorFabricErrorServiceTypeNotFoundFabricErrorServiceTypeMismatchFabricErrorServiceTypeTemplateNotFoundFabricErrorConfigurationSectionNotFoundFabricErrorConfigurationParameterNotFoundFabricErrorInvalidConfigurationFabricErrorImagebuilderValidationErrorFabricErrorPartitionNotFoundFabricErrorReplicaDoesNotExistFabricErrorServiceGroupAlreadyExistsFabricErrorServiceGroupDoesNotExistFabricErrorProcessDeactivatedFabricErrorProcessAbortedFabricErrorUpgradeFailedFabricErrorInvalidCredentialTypeFabricErrorInvalidX509FindTypeFabricErrorInvalidX509StoreLocationFabricErrorInvalidX509StoreNameFabricErrorInvalidX509ThumbprintFabricErrorInvalidProtectionLevelFabricErrorInvalidX509StoreFabricErrorInvalidSubjectNameFabricErrorInvalidAllowedCommonNameListFabricErrorInvalidCredentialsFabricErrorDecryptionFailedFabricErrorConfigurationPackageNotFoundFabricErrorDataPackageNotFoundFabricErrorCodePackageNotFoundFabricErrorServiceEndpointResourceNotFoundFabricErrorInvalidOperationFabricErrorObjectClosedFabricErrorTimeoutFabricErrorFileNotFoundFabricErrorDirectoryNotFoundFabricErrorInvalidDirectoryFabricErrorPathTooLongFabricErrorImagestoreIoerrorFabricErrorCorruptedImageStoreObjectFoundFabricErrorApplicationNotUpgradingFabricErrorApplicationAlreadyInTargetVersionFabricErrorImagebuilderUnexpectedErrorFabricErrorFabricVersionNotFoundFabricErrorFabricVersionInUseFabricErrorFabricVersionAlreadyExistsFabricErrorFabricAlreadyInTargetVersionFabricErrorFabricNotUpgradingFabricErrorFabricUpgradeInProgressFabricErrorFabricUpgradeValidationErrorFabricErrorHealthMaxReportsReachedFabricErrorHealthStaleReportFabricErrorKeyTooLargeFabricErrorKeyNotFoundFabricErrorSequenceNumberCheckFailedFabricErrorEncryptionFailedFabricErrorInvalidAtomicGroupFabricErrorHealthEntityNotFoundFabricErrorServiceManifestNotFoundFabricErrorReliableSessionTransportStartupFailureFabricErrorReliableSessionAlreadyExistsFabricErrorReliableSessionCannotConnectFabricErrorReliableSessionManagerExistsFabricErrorReliableSessionRejectedFabricErrorReliableSessionManagerAlreadyListeningFabricErrorReliableSessionManagerNotFoundFabricErrorReliableSessionManagerNotListeningFabricErrorInvalidServiceTypeFabricErrorImagebuilderTimeoutFabricErrorImagebuilderAccessDeniedFabricErrorImagebuilderInvalidMsiFileFabricErrorServiceTooBusyFabricErrorTransactionNotActiveFabricErrorRepairTaskAlreadyExistsFabricErrorRepairTaskNotFoundFabricErrorReliableSessionNotFoundFabricErrorReliableSessionQueueEmptyFabricErrorReliableSessionQuotaExceededFabricErrorReliableSessionServiceFaultedFabricErrorReliableSessionInvalidTargetPartitionFabricErrorTransactionTooLargeFabricErrorReplicationOperationTooLargeFabricErrorInstanceIdMismatchFabricErrorUpgradeDomainAlreadyCompletedFabricErrorNodeHasNotStoppedYetFabricErrorInsufficientClusterCapacityFabricErrorInvalidPackageSharingPolicyFabricErrorPredeploymentNotAllowedFabricErrorInvalidBackupSettingFabricErrorMissingFullBackupFabricErrorBackupInProgressFabricErrorDuplicateServiceNotificationFilterNameFabricErrorInvalidReplicaOperationFabricErrorInvalidReplicaStateFabricErrorLoadbalancerNotReadyFabricErrorInvalidPartitionOperationFabricErrorPrimaryAlreadyExistsFabricErrorSecondaryAlreadyExistsFabricErrorBackupDirectoryNotEmptyFabricErrorForceNotSupportedForReplicaOperationFabricErrorAcquireFileLockFailedFabricErrorConnectionDeniedFabricErrorServerAuthenticationFailedFabricErrorConstraintKeyUndefinedFabricErrorMultithreadedTransactionsNotAllowedFabricErrorInvalidX509NameListFabricErrorVerboseFmPlacementHealthReportingRequiredFabricErrorGatewayNotReachableFabricErrorUserRoleClientCertificateNotConfiguredFabricErrorTransactionAbortedFabricErrorCannotConnectFabricErrorMessageTooLargeFabricErrorConstraintNotSatisfiedFabricErrorEndpointNotFoundFabricErrorApplicationUpdateInProgressFabricErrorDeleteBackupFileFailedFabricErrorConnectionClosedByRemoteEndFabricErrorInvalidTestCommandStateFabricErrorTestCommandOperationIdAlreadyExistsFabricErrorCmOperationFailedFabricErrorImagebuilderReservedDirectoryErrorFabricErrorCertificateNotFoundFabricErrorChaosAlreadyRunningFabricErrorFabricDataRootNotFoundFabricErrorInvalidRestoreDataFabricErrorDuplicateBackupsFabricErrorInvalidBackupChainFabricErrorStopInProgressFabricErrorAlreadyStoppedFabricErrorNodeIsDownFabricErrorNodeTransitionInProgressFabricErrorInvalidBackupFabricErrorInvalidInstanceIdFabricErrorInvalidDurationFabricErrorRestoreSafeCheckFailedFabricErrorConfigUpgradeFailedFabricErrorUploadSessionRangeNotSatisfiableFabricErrorUploadSessionIdConflictFabricErrorInvalidPartitionSelectorFabricErrorInvalidReplicaSelectorFabricErrorDnsServiceNotFoundFabricErrorInvalidDnsNameFabricErrorDnsNameInUseFabricErrorComposeDeploymentAlreadyExistsFabricErrorComposeDeploymentNotFoundFabricErrorInvalidForStatefulServicesFabricErrorInvalidForStatelessServicesFabricErrorOnlyValidForStatefulPersistentServicesFabricErrorInvalidUploadSessionIdFabricErrorBackupNotEnabledFabricErrorBackupIsEnabledFabricErrorBackupPolicyDoesNotExistFabricErrorBackupPolicyAlreadyExistsFabricErrorRestoreInProgressFabricErrorRestoreSourceTargetPartitionMismatchFabricErrorFaultAnalysisServiceNotEnabledFabricErrorContainerNotFoundFabricErrorObjectDisposedFabricErrorNotReadableFabricErrorBackupcopierUnexpectedErrorFabricErrorBackupcopierTimeoutFabricErrorBackupcopierAccessDeniedFabricErrorInvalidServiceScalingPolicyFabricErrorSingleInstanceApplicationAlreadyExistsFabricErrorSingleInstanceApplicationNotFoundFabricErrorVolumeAlreadyExistsFabricErrorVolumeNotFoundFabricErrorDatabaseMigrationInProgressFabricErrorCentralSecretServiceGenericFabricErrorSecretInvalidFabricErrorSecretVersionAlreadyExistsFabricErrorSingleInstanceApplicationUpgradeInProgressFabricErrorOperationNotSupportedFabricErrorComposeDeploymentNotUpgradingFabricErrorSecretTypeCannotBeChangedFabricErrorNetworkNotFoundFabricErrorNetworkInUseFabricErrorEndpointNotReferenced"
	_FabricErrorCode_name_1 = "FabricErrorLastReservedHresult"
)

var (
	_FabricErrorCode_index_0 = [...]uint16{0, 31, 56, 81, 111, 139, 166, 189, 212, 231, 255, 276, 295, 326, 357, 390, 421, 452, 482, 507, 541, 584, 623, 658, 682, 703, 733, 757, 788, 833, 872, 906, 937, 972, 1002, 1041, 1085, 1115, 1145, 1183, 1222, 1263, 1294, 1332, 1360, 1390, 1426, 1461, 1490, 1515, 1539, 1571, 1601, 1636, 1667, 1699, 1732, 1759, 1788, 1827, 1856, 1883, 1922, 1952, 1982, 2024, 2051, 2074, 2092, 2115, 2143, 2170, 2192, 2220, 2261, 2295, 2339, 2377, 2409, 2438, 2475, 2514, 2543, 2577, 2616, 2650, 2678, 2700, 2722, 2758, 2785, 2814, 2845, 2879, 2928, 2967, 3006, 3045, 3079, 3128, 3169, 3214, 3243, 3273, 3308, 3345, 3370, 3401, 3435, 3464, 3498, 3534, 3573, 3613, 3661, 3691, 3730, 3759, 3799, 3830, 3868, 3906, 3940, 3971, 3999, 4026, 4075, 4109, 4139, 4170, 4206, 4237, 4270, 4304, 4351, 4383, 4410, 4447, 4480, 4526, 4556, 4608, 4638, 4687, 4716, 4740, 4766, 4799, 4826, 4864, 4897, 4935, 4969, 5015, 5043, 5088, 5118, 5148, 5181, 5210, 5237, 5266, 5291, 5316, 5337, 5372, 5396, 5424, 5450, 5483, 5513, 5556, 5590, 5625, 5658, 5687, 5712, 5735, 5776, 5812, 5849, 5887, 5936, 5969, 5996, 6022, 6057, 6093, 6121, 6168, 6209, 6237, 6262, 6284, 6322, 6352, 6387, 6425, 6474, 6518, 6548, 6573, 6611, 6649, 6673, 6710, 6763, 6795, 6835, 6871, 6897, 6920, 6952}
)

func (i FabricErrorCode) String() string {
	switch {
	case 2147949500 <= i && i <= 2147949711:
		i -= 2147949500
		return _FabricErrorCode_name_0[_FabricErrorCode_index_0[i]:_FabricErrorCode_index_0[i+1]]
	case i == 2147949899:
		return _FabricErrorCode_name_1
	default:
		return "FabricErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
