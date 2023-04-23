package dashboardbuilder

import (
	"fmt"
	"os"
	"time"

	"github.com/mum4k/termdash"

	domainhost "github.com/bartosian/sui_helpers/suimon/internal/core/domain/host"
	"github.com/bartosian/sui_helpers/suimon/internal/core/domain/service/dashboardbuilder/dashboards"
)

const (
	renderInterval = 500 * time.Millisecond
	queryInterval  = 1 * time.Second
)

// Render displays the dashboard on the terminal and updates the cells with new data periodically.
func (db *Builder) Render() (err error) {
	// Use a deferred function to call db.TearDown() if there were errors or panics
	defer func() {
		if err != nil {
			db.tearDown()
		}

		if err := recover(); err != nil {
			// Handle the panic by logging the error and exiting the program
			db.tearDown()

			db.cliGateway.Error(fmt.Sprintf("panic: %v", err))

			os.Exit(1)
		}
	}()

	// Create a channel for error handling
	errChan := make(chan error, 10)
	defer close(errChan)

	tickerQuery := time.NewTicker(queryInterval)
	defer tickerQuery.Stop()

	// Start a goroutine for the metric retrieval loop
	go func(host *domainhost.Host) {
		for {
			select {
			case <-tickerQuery.C:
				if err := host.GetMetrics(); err != nil {
					errChan <- err

					return
				}
			case <-db.ctx.Done():
				return
			}
		}
	}(&db.host)

	tickerRerender := time.NewTicker(renderInterval)
	defer tickerRerender.Stop()

	// Start a goroutine for the dashboard rendering loop
	go func(host *domainhost.Host) {
		for {
			select {
			case <-tickerRerender.C:
				columns := dashboards.GetNodeColumnValues(host)

				for columnName, cell := range db.cells {
					columnValue, ok := columns[columnName]
					if !ok {
						errChan <- fmt.Errorf("failed to get metric for column %s", columnName)

						return
					}

					// Write the new data to the cell
					if err := cell.Write(columnValue); err != nil {
						errChan <- err

						return
					}
				}
			case <-db.ctx.Done():
				return
			}
		}
	}(&db.host)

	// Display the dashboard on the terminal and handle errors
	if err := termdash.Run(
		db.ctx, db.terminal, db.dashboard,
		termdash.KeyboardSubscriber(db.quitter),
		termdash.RedrawInterval(renderInterval),
	); err != nil {
		return fmt.Errorf("failed to run terminal dashboard: %w", err)
	}

	// Check for errors from the ticker loop
	select {
	case err := <-errChan:
		return err
	case <-db.ctx.Done():
		return nil
	}
}
