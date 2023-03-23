package columnnames

type SystemColumnName int

const (
	SystemColumnNameEpoch SystemColumnName = iota
	SystemColumnNameEpochDurationMs
	SystemColumnNameStorageFund
	SystemColumnNameReferenceGasPrice
	SystemColumnNameStakeSubsidyCounter
	SystemColumnNameStakeSubsidyBalance
	SystemColumnNameStakeSubsidyCurrentEpochAmount
	SystemColumnNameTotalStake
	SystemColumnNameValidatorsCount
	SystemColumnNamePendingActiveValidatorsSize
	SystemColumnNamePendingRemovals
	SystemColumnNameValidatorsCandidateSize
	SystemColumnNameValidatorsAtRiskCount
)
