package main

import (
	"fmt"
	"os"

	"github.com/ledgerwatch/erigon/common"
	"github.com/ledgerwatch/erigon/params"
	"github.com/spf13/cobra"
)

var (
	contractCodePath string
	calldata         string
	numRuns          int
)

func check(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}

var cmd = &cobra.Command{
	Use:   "runner-erigon",
	Short: "erigon runner for evm-bench",
	Run: func(_ *cobra.Command, _ []string) {
		contractCodeHex, err := os.ReadFile(contractCodePath)
		check(err)

		contractCodeBytes := common.Hex2Bytes(string(contractCodeHex))
		calldataBytes := common.Hex2Bytes(calldata)

		// fmt.Println(contractCodeBytes)
		// fmt.Println(calldataBytes)

		zeroAddress := libcommon.BytesToAddress(common.FromHex("0x0000000000000000000000000000000000000000"))
		callerAddress := libcommon.BytesToAddress(common.FromHex("0x1000000000000000000000000000000000000001"))

		// fmt.Println(zeroAddress)
		// fmt.Println(callerAddress)
		var config = params.MainnetChainConfig
		rules := config.Rules(uint64(12965000), uint64(1681338458))
		fmt.Println(config)
		fmt.Println(rules)

		// genesis := &core.Genesis{
		// 	Config:     config,
		// 	Coinbase:   zeroAddress,
		// 	Difficulty: Difficulty: big.NewInt(17179869184),,
		// 	GasLimit:   5000,
		// 	Number:     1681338457,
		// 	Timestamp:  1681338458,
		// 	Alloc:      core.readPrealloc("alloc/mainnet.json"),
		// }

		// fmt.Println(genesis.Alloc)
		// for i := 0; i < numRuns; i++ {
		// 	fmt.Println("..")
		// }

	},
}

func init() {
	cmd.Flags().StringVar(&contractCodePath, "contract-code-path", "", "Path to the hex contract code to deploy and run")
	cmd.MarkFlagRequired("contract-code-path")
	cmd.Flags().StringVar(&calldata, "calldata", "", "Hex of calldata to use when calling the contract")
	cmd.MarkFlagRequired("calldata")
	cmd.Flags().IntVar(&numRuns, "num-runs", 0, "Number of times to run the benchmark")
	cmd.MarkFlagRequired("num-runs")
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
