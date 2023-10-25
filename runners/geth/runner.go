package main

import (
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
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
	Use:   "runner-geth",
	Short: "go-ethereum runner for evm-bench",
	Run: func(_ *cobra.Command, _ []string) {
		contractCodeHex, err := os.ReadFile(contractCodePath)
		check(err)

		contractCodeBytes := common.Hex2Bytes(string(contractCodeHex))
		calldataBytes := common.Hex2Bytes(calldata)

		zeroAddress := common.BytesToAddress(common.FromHex("0x0000000000000000000000000000000000000000"))
		callerAddress := common.BytesToAddress(common.FromHex("0x1000000000000000000000000000000000000001"))

		config := params.MainnetChainConfig
		rules := config.Rules(config.LondonBlock, true, 1681338458)
		defaultGenesis := core.DefaultGenesisBlock()
		genesis := &core.Genesis{
			Config:     config,
			Coinbase:   defaultGenesis.Coinbase,
			Difficulty: defaultGenesis.Difficulty,
			GasLimit:   defaultGenesis.GasLimit,
			Number:     1681338457,
			Timestamp:  1681338458,
			Alloc:      defaultGenesis.Alloc,
		}

		statedb, err := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
		check(err)
		//fmt.Fprintln(os.Stdout, "hahaha")
		//fmt.Println("ok!1111")

		zeroValue := big.NewInt(0)
		gasLimit := ^uint64(0)

		tx := types.NewTx(&types.AccessListTx{
			ChainID:  big.NewInt(1),
			Nonce:    0,
			To:       &zeroAddress,
			Value:    zeroValue,
			Gas:      gasLimit,
			GasPrice: zeroValue,
			Data:     contractCodeBytes,
		})

		var signer types.Signer
		signer = types.NewEIP2930Signer(big.NewInt(1))
		signer.Sender(tx)
		createMsg, _ := core.TransactionToMessage(tx, signer, zeroValue)

		//createMsg := types.NewMessage(callerAddress, &zeroAddress, 0, zeroValue, gasLimit, zeroValue, zeroValue, zeroValue, contractCodeBytes, types.AccessList{}, false)
		statedb.Prepare(rules, callerAddress, zeroAddress, &zeroAddress, vm.ActivePrecompiles(rules), createMsg.AccessList)

		blockContext := core.NewEVMBlockContext(genesis.ToBlock().Header(), nil, &zeroAddress)
		txContext := core.NewEVMTxContext(createMsg)
		evm := vm.NewEVM(blockContext, txContext, statedb, config, vm.Config{})

		fmt.Println("chainconfig:  ", evm.ChainConfig())
		_, contractAddress, _, err := evm.Create(vm.AccountRef(callerAddress), contractCodeBytes, gasLimit, new(big.Int))
		check(err)

		//fmt.Fprintln(os.Stdout, "hehehe")
		//fmt.Println("ok!2222")

		tx1 := types.NewTx(&types.AccessListTx{
			ChainID:  big.NewInt(1),
			Nonce:    1,
			To:       &contractAddress,
			Value:    zeroValue,
			Gas:      gasLimit,
			GasPrice: zeroValue,
			Data:     calldataBytes,
		})

		signer.Sender(tx1)
		msg, _ := core.TransactionToMessage(tx1, signer, zeroValue)
		//msg := types.NewMessage(callerAddress, &contractAddress, 1, zeroValue, gasLimit, zeroValue, zeroValue, zeroValue, calldataBytes, types.AccessList{}, false)

		for i := 0; i < numRuns; i++ {
			snapshot := statedb.Snapshot()
			statedb.AddAddressToAccessList(msg.From)
			statedb.AddAddressToAccessList(*msg.To)

			start := time.Now()
			_, _, err := evm.Call(vm.AccountRef(callerAddress), *msg.To, msg.Data, msg.GasLimit, msg.Value)
			timeTaken := time.Since(start)

			fmt.Println(float64(timeTaken.Microseconds()) / 1e3)

			check(err)

			statedb.RevertToSnapshot(snapshot)
		}
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
