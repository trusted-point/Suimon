// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"

	"github.com/bartosian/sui_helpers/suimon/internal/core/controllers"
	"github.com/bartosian/sui_helpers/suimon/internal/core/controllers/monitor"
	"github.com/bartosian/sui_helpers/suimon/internal/core/domain/config"
	"github.com/bartosian/sui_helpers/suimon/internal/core/gateways/cligw"
	"github.com/bartosian/sui_helpers/suimon/internal/core/gateways/geogw"
	"github.com/bartosian/sui_helpers/suimon/internal/core/gateways/prometheusgw"
	"github.com/bartosian/sui_helpers/suimon/internal/core/gateways/rpcgw"
	"github.com/bartosian/sui_helpers/suimon/internal/core/handlers/commands"
	"github.com/bartosian/sui_helpers/suimon/internal/pkg/log"
)

func main() {
	var logger = log.NewLogger()

	cliGateway := cligw.NewCliGateway()

	config, err := config.NewConfig(logger)
	if err != nil {
		// If an error occurs during initialization of the tables object, log the error and exit the program.
		cliGateway.Error(err.Error())
		return
	}

	// Instantiate gateways
	rpcGateway := rpcgw.NewGateway(logger, "")
	geoGateway := geogw.NewGateway(logger, "", "")
	prometheusGateway := prometheusgw.NewGateway(logger, "")

	// Instantiate controllers
	rootController := controllers.NewRootController(cliGateway)
	versionController := controllers.NewVersionController(cliGateway)
	monitorController := monitor.NewController(
		logger,
		config,
		rpcGateway,
		geoGateway,
		prometheusGateway,
		cliGateway,
	)

	// Instantiate Handlers - Root
	rootCmdHandler := cmdhandlers.NewRootHandler(rootController)

	// Instantiate Handlers - second level
	versionCmdHandler := cmdhandlers.NewVersionHandler(versionController)
	monitorCmdHandler := cmdhandlers.NewMonitorHandler(monitorController)

	// Add subcommands to the root command handler
	rootCmdHandler.AddSubCommands(
		versionCmdHandler,
		monitorCmdHandler,
	)

	// Start the root command handler
	rootCmdHandler.Start()
}

func handlePanic(logger *log.Logger) {
	if r := recover(); r != nil {
		logger.Error("failed to execute suimon, please check an issue: ", r)

		os.Exit(1)
	}
}
