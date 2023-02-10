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
	"flag"
	"time"

	"github.com/bartosian/sui_helpers/suimon/cmd/checker"
	"github.com/bartosian/sui_helpers/suimon/cmd/checker/config"
	"github.com/bartosian/sui_helpers/suimon/pkg/log"

	"github.com/schollz/progressbar/v3"
)

var (
	suimonConfigPath = flag.String("sf", "", "(optional) path to the suimon config file, can use SUIMON_CONFIG_PATH env variable instead")
	nodeConfigPath   = flag.String("nf", "", "(optional) path to the node config file, can use SUIMON_NODE_CONFIG_PATH variable instead")
	network          = flag.String("n", "", "(optional) network name, possible values: testnet, devnet")
)

const (
	suimonConfigNotFound = `provide path to the suimon.yaml file by using -sf option 
or by setting SUIMON_CONFIG_PATH env variable
or put suimon.yaml in $HOME/.suimon/suimon.yaml`
	nodeConfigNotFound = `provide path to the fullnode.yaml file by using -sf option
or by setting SUIMON_NODE_CONFIG_PATH env variable
or set path to this file in suimon.yaml`
	invalidNetworkTypeProvided = `provide valid network type by using -n option
or set it in suimon.yaml`
)

func main() {
	flag.Parse()
	logger := log.NewLogger()

	progressChan := make(chan struct{})

	// start showing progress bar
	go newProgressBar(progressChan)

	// parse suimon.yaml config file
	suimonConfig, err := config.ParseSuimonConfig(suimonConfigPath)
	if err != nil {
		logger.Error(suimonConfigNotFound)

		return
	}

	// parse fullnode/validator.yaml config file
	nodeConfig, err := config.ParseNodeConfig(nodeConfigPath, suimonConfig.NodeConfigPath)
	if err != nil {
		logger.Error(nodeConfigNotFound)

		return
	}

	// parse network flag
	networkConfig, err := config.ParseNetworkConfig(suimonConfig, network)
	if err != nil {
		logger.Error(invalidNetworkTypeProvided)

		return
	}

	suimonConfig.SetNetworkConfig(networkConfig)

	// create checker instance to process to request all the required data and pass them to tablebuilder
	checker, err := checker.NewChecker(*suimonConfig, *nodeConfig)
	if err != nil {
		logger.Error("failed to create peers checker: ", err)

		return
	}

	if err := checker.ParseData(); err != nil {
		logger.Error("failed to parse data: ", err)

		return
	}

	// initialize tables with the styles and data received
	checker.GenerateSystemTable()
	checker.GenerateNodeTable()
	checker.GeneratePeersTable()

	// stop showing progress bar
	progressChan <- struct{}{}
	
	// draw initialized tables to the terminal
	checker.DrawTable()
}

func newProgressBar(progressChan chan struct{}) {
	progressTicker := time.NewTicker(20 * time.Millisecond)

	bar := progressbar.NewOptions(1000,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionSetWidth(25),
		progressbar.OptionSetDescription(" [ GENERATING TABLES... ] "),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	for {
		select {
		case <-progressChan:
			progressTicker.Stop()
			bar.Clear()

			return
		case <-progressTicker.C:
			for i := 0; i < 500; i++ {
				bar.Add(1)

				time.Sleep(7 * time.Millisecond)
			}
		}
	}
}
