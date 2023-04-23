package dashboards

import (
	"fmt"

	"github.com/bartosian/sui_helpers/suimon/internal/core/domain/enums"
	"github.com/bartosian/sui_helpers/suimon/internal/core/domain/host"
)

var (
	ColumnsConfigNode = ColumnsConfig{
		enums.ColumnNameCurrentEpoch:                 19,
		enums.ColumnNameNetworkPeers:                 19,
		enums.ColumnNameUptime:                       19,
		enums.ColumnNameVersion:                      19,
		enums.ColumnNameCommit:                       19,
		enums.ColumnNameHealth:                       5,
		enums.ColumnNameTotalTransactionBlocks:       33,
		enums.ColumnNameTotalTransactionCertificates: 33,
		enums.ColumnNameTotalTransactionEffects:      33,
		enums.ColumnNameLatestCheckpoint:             24,
		enums.ColumnNameHighestKnownCheckpoint:       24,
		enums.ColumnNameHighestSyncedCheckpoint:      24,
		enums.ColumnNameLastExecutedCheckpoint:       24,
		enums.ColumnNameCheckpointExecBacklog:        24,
		enums.ColumnNameCheckpointSyncBacklog:        24,
		enums.ColumnNameTXSyncPercentage:             24,
		enums.ColumnNameCheckSyncPercentage:          24,
	}

	RowsConfigNode = RowsConfig{
		0: {
			Height: 15,
			Columns: []enums.ColumnName{
				enums.ColumnNameCurrentEpoch,
				enums.ColumnNameNetworkPeers,
				enums.ColumnNameUptime,
				enums.ColumnNameVersion,
				enums.ColumnNameCommit,
			},
		},
		1: {
			Height: 15,
			Columns: []enums.ColumnName{
				enums.ColumnNameTotalTransactionBlocks,
				enums.ColumnNameTotalTransactionCertificates,
				enums.ColumnNameTotalTransactionEffects,
			},
		},
		2: {
			Height: 15,
			Columns: []enums.ColumnName{
				enums.ColumnNameLatestCheckpoint,
				enums.ColumnNameHighestKnownCheckpoint,
				enums.ColumnNameHighestSyncedCheckpoint,
				enums.ColumnNameLastExecutedCheckpoint,
			},
		},
		3: {
			Height: 15,
			Columns: []enums.ColumnName{
				enums.ColumnNameTXSyncPercentage,
				enums.ColumnNameCheckSyncPercentage,
				enums.ColumnNameCheckpointExecBacklog,
				enums.ColumnNameCheckpointSyncBacklog,
			},
		},
	}

	CellsConfigNode = CellsConfig{
		enums.ColumnNameHealth:                       "HEALTH",
		enums.ColumnNameTotalTransactionBlocks:       "TOTAL TRANSACTION BLOCKS",
		enums.ColumnNameLatestCheckpoint:             "LATEST CHECKPOINT",
		enums.ColumnNameTotalTransactionCertificates: "TOTAL TRANSACTION CERTIFICATES",
		enums.ColumnNameTotalTransactionEffects:      "TOTAL TRANSACTION EFFECTS",
		enums.ColumnNameHighestKnownCheckpoint:       "HIGHEST KNOWN CHECKPOINT",
		enums.ColumnNameHighestSyncedCheckpoint:      "HIGHEST SYNCED CHECKPOINT",
		enums.ColumnNameLastExecutedCheckpoint:       "LAST EXECUTED CHECKPOINT",
		enums.ColumnNameCheckpointExecBacklog:        "CHECKPOINT EXEC BACKLOG",
		enums.ColumnNameCheckpointSyncBacklog:        "CHECKPOINT SYNC BACKLOG",
		enums.ColumnNameCurrentEpoch:                 "CURRENT EPOCH",
		enums.ColumnNameTXSyncPercentage:             "TX SYNC PERCENTAGE",
		enums.ColumnNameCheckSyncPercentage:          "CHECKPOINTS SYNC PERCENTAGE",
		enums.ColumnNameNetworkPeers:                 "NETWORK PEERS",
		enums.ColumnNameUptime:                       "UPTIME",
		enums.ColumnNameVersion:                      "VERSION",
		enums.ColumnNameCommit:                       "COMMIT",
	}
)

// GetNodeColumnValues returns a map of NodeColumnName values to corresponding values for a node at the specified index on the specified host.
// The function retrieves information about the node from the host's internal state and formats it into a map of NodeColumnName keys and corresponding values.
// The function also includes emoji values in the map if the specified flag is true.
// Returns a map of NodeColumnName keys to corresponding values.
func GetNodeColumnValues(host *host.Host) ColumnValues {
	status := host.Status.StatusToPlaceholder()

	columnValues := ColumnValues{
		enums.ColumnNameHealth:                       status,
		enums.ColumnNameTotalTransactionBlocks:       host.Metrics.TotalTransactionsBlocks,
		enums.ColumnNameTotalTransactionCertificates: host.Metrics.TotalTransactionCertificates,
		enums.ColumnNameTotalTransactionEffects:      host.Metrics.TotalTransactionEffects,
		enums.ColumnNameLatestCheckpoint:             host.Metrics.LatestCheckpoint,
		enums.ColumnNameHighestKnownCheckpoint:       host.Metrics.HighestKnownCheckpoint,
		enums.ColumnNameHighestSyncedCheckpoint:      host.Metrics.HighestSyncedCheckpoint,
		enums.ColumnNameLastExecutedCheckpoint:       host.Metrics.LastExecutedCheckpoint,
		enums.ColumnNameCheckpointExecBacklog:        host.Metrics.CheckpointExecBacklog,
		enums.ColumnNameCheckpointSyncBacklog:        host.Metrics.CheckpointSyncBacklog,
		enums.ColumnNameCurrentEpoch:                 host.Metrics.CurrentEpoch,
		enums.ColumnNameTXSyncPercentage:             fmt.Sprintf("%v%%", host.Metrics.TxSyncPercentage),
		enums.ColumnNameCheckSyncPercentage:          fmt.Sprintf("%v%%", host.Metrics.CheckSyncPercentage),
		enums.ColumnNameNetworkPeers:                 host.Metrics.NetworkPeers,
		enums.ColumnNameUptime:                       host.Metrics.Uptime,
		enums.ColumnNameVersion:                      host.Metrics.Version,
		enums.ColumnNameCommit:                       host.Metrics.Commit,
	}

	return columnValues
}
