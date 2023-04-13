package tablebuilder

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/bartosian/sui_helpers/suimon/internal/core/domain/enums"
	"github.com/bartosian/sui_helpers/suimon/internal/core/domain/host"
	"github.com/bartosian/sui_helpers/suimon/internal/core/domain/service/tablebuilder/tables"
	"github.com/bartosian/sui_helpers/suimon/internal/core/gateways/cligw"
)

type (
	Builders struct {
		peerTable             Builder
		nodeTable             Builder
		validatorTable        Builder
		rpcTable              Builder
		systemStateTable      Builder
		validatorCountsTable  Builder
		atRiskValidatorsTable Builder
		validatorReportsTable Builder
		activeValidatorsTable Builder
	}

	Builder struct {
		tableType  enums.TableType
		hosts      []host.Host
		cliGateway *cligw.Gateway
		writer     table.Writer
		config     *tables.TableConfig
	}
)

// NewBuilder creates a new instance of the table builder, using the CLI gateway
func NewBuilder(tableType enums.TableType, hosts []host.Host, cliGateway *cligw.Gateway) *Builder {
	tableWR := table.NewWriter()
	tableWR.SetOutputMirror(os.Stdout)

	return &Builder{
		tableType:  tableType,
		hosts:      hosts,
		cliGateway: cliGateway,
		writer:     tableWR,
	}
}

// setColumns sets the column configurations for the table builder based on the configuration in the builder's table config
func (tb *Builder) setColumns() {
	var columnsConfig []table.ColumnConfig

	for _, column := range tb.config.Columns {
		columnsConfig = append(columnsConfig, *column.Config)
	}

	tb.writer.SetColumnConfigs(columnsConfig)
}

// setRows sets the rows of the table builder based on the configuration in the builder's table config
func (tb *Builder) setRows() {
	rowsConfig := tb.config.Rows
	columnsConfig := tb.config.Columns
	itemsCount := tb.config.RowsCount
	columnsPerRow := len(rowsConfig[0])

	for itemIndex := 0; itemIndex < itemsCount; itemIndex++ {
		for rowIndex, columns := range rowsConfig {
			header := tables.NewRow(true, false, columnsPerRow, true, text.AlignCenter)
			footer := tables.NewRow(false, false, columnsPerRow, true, text.AlignCenter)
			row := tables.NewRow(false, true, columnsPerRow, true, text.AlignCenter)

			var (
				columnIdx  int
				columnName enums.ColumnName
			)

			for columnIdx, columnName = range columns {
				columnConfig := columnsConfig[columnName]
				columnValue := columnConfig.Values[itemIndex]

				header.AppendValue(columnName.ToString())
				row.AppendValue(columnValue)
				footer.PrependValue(tables.EmptyValue)
			}

			columnIdx++

			for columnIdx < columnsPerRow {
				header.PrependValue(tables.EmptyValue)
				footer.PrependValue(tables.EmptyValue)
				row.PrependValue(tables.EmptyValue)

				columnIdx++
			}

			if itemIndex == 0 && rowIndex == 0 {
				tb.writer.AppendHeader(header.Values, header.Config)
				tb.writer.AppendFooter(footer.Values, footer.Config)
			} else if rowIndex%2 == 1 || itemIndex > 0 && len(rowsConfig) > 1 && rowIndex%2 == 0 {
				tb.writer.AppendRow(header.Values, header.Config)
			}

			tb.writer.AppendRow(row.Values, row.Config)
			tb.writer.AppendSeparator()
		}
	}
}

// setStyle sets the style for the table builder based on the configuration in the builder's table config
func (tb *Builder) setStyle() {
	tb.writer.SortBy(tb.config.Sort)
	tb.writer.SetTitle(tb.config.Name)
	tb.writer.SetStyle(tb.config.Style)
	tb.writer.SetAutoIndex(tb.config.AutoIndex)

	tb.setColors()
}

// setColors sets the row colors for the table builder based on the current state of the table
func (tb *Builder) setColors() {
	var f = func() func(row table.Row) text.Colors {
		valuesRowFgColor := text.Colors{text.FgWhite}
		bgColor := []text.Color{text.BgWhite, text.BgHiBlue, text.BgHiBlue, text.BgWhite}
		currentColor := 0

		var handler = func(row table.Row) text.Colors {
			for _, column := range row {
				if _, ok := column.(int); ok {
					return valuesRowFgColor
				}
			}

			colors := text.Colors{text.FgBlack, bgColor[currentColor]}

			currentColor++
			if currentColor > 3 {
				currentColor = 0
			}

			return colors
		}

		return handler
	}()

	tb.writer.SetRowPainter(f)
}
