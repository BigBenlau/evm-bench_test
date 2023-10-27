package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/holiman/uint256"
	libcommon "github.com/ledgerwatch/erigon-lib/common"
	"github.com/ledgerwatch/erigon-lib/common/hexutility"
	"github.com/ledgerwatch/erigon-lib/kv/memdb"
	"github.com/ledgerwatch/erigon/common"
	"github.com/ledgerwatch/erigon/core"
	"github.com/ledgerwatch/erigon/core/state"
	"github.com/ledgerwatch/erigon/core/types"
	"github.com/ledgerwatch/erigon/core/vm"
	"github.com/ledgerwatch/erigon/crypto"
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

		contractCodeBytes := hexutility.Hex2Bytes(string(contractCodeHex))
		calldataBytes := hexutility.Hex2Bytes(calldata)
		// fmt.Println(contractCodeBytes)
		// fmt.Println(calldataBytes)

		zerodata := uint256.NewInt(0)
		zeroAddress := libcommon.BytesToAddress(common.FromHex("0x0000000000000000000000000000000000000000"))
		callerAddress := libcommon.BytesToAddress(common.FromHex("0x1000000000000000000000000000000000000001"))
		// callerAddress := libcommon.BytesToAddress(common.FromHex("0x1000000000000000000000000000000000000001"))
		// fmt.Println(zeroAddress)
		// fmt.Println(callerAddress)

		var config = params.MainnetChainConfig
		rules := config.Rules(uint64(12965000), uint64(1681338458))
		// fmt.Println(config)
		// fmt.Println(rules)
		gasLimit := ^uint64(5000)

		//create genesis
		var alloc types.GenesisAlloc

		// address1 := libcommon.HexToAddress("0x000d836201318ec6899a67540690382780743280")
		// address2 := libcommon.HexToAddress("0x001762430ea9c3a26e5749afdb70da5f78ddbb8c")
		// funds := new(big.Int)
		// funds.SetString("ad78ebc5ac6200000", 16)

		// alloc = types.GenesisAlloc{
		// 	address1:      {Balance: funds},
		// 	address2:      {Balance: funds},
		// 	callerAddress: {Balance: funds},
		// }

		alloc = core.MainnetGenesisBlock().Alloc
		funds := new(big.Int)
		funds.SetString("ad78ebc5ac62000", 16)
		addalloc := types.GenesisAccount{
			Balance: funds,
		}
		alloc[callerAddress] = addalloc

		genesis := &types.Genesis{
			Config:     config,
			Coinbase:   zeroAddress,
			Difficulty: big.NewInt(17179869184),
			GasLimit:   core.MainnetGenesisBlock().GasLimit,
			Number:     1681338457,
			Timestamp:  1681338458,
			Alloc:      alloc,
		}

		block, statedb, err := core.GenesisToBlock(genesis, "")
		db := memdb.New("")
		ctx, err := db.BeginRw(context.Background())
		statedb = state.New(state.NewPlainStateReader(ctx))
		//m := mock.MockWithGenesisEngine(t, bt.genesis(config), engine, false, checkStateRoot)
		//NewPlainStateReader(tx)
		//statedb, err := state.New()

		check(err)

		//create a new trans and mess
		key, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(key.PublicKey)
		signer := types.LatestSignerForChainID(big.NewInt(18))
		statedb.AddAddressToAccessList(addr)
		tx, err := types.SignTx(types.NewTransaction(0, libcommon.Address{}, zerodata, gasLimit, zerodata, contractCodeBytes), *signer, key)
		fmt.Println("tx:  ", tx)
		createMsg, _ := tx.AsMessage(*signer, big.NewInt(0), rules)

		//prepare statedb
		accesslist := createMsg.AccessList()
		statedb.Prepare(rules, callerAddress, zeroAddress, &zeroAddress, vm.ActivePrecompiles(rules), accesslist)

		//create blockContext and txContext
		blockContext := core.NewEVMBlockContext(block.Header(), core.GetHashFn(block.Header(), nil), nil, &zeroAddress)
		txContext := core.NewEVMTxContext(createMsg)
		fmt.Println("blockContext:  ", blockContext)

		//new a evm
		evm := vm.NewEVM(blockContext, txContext, statedb, config, vm.Config{})
		fmt.Println("evm:  ", evm)
		fmt.Println("chainconfig:  ", evm.ChainConfig())

		//create contract
		var value *uint256.Int
		value = uint256.NewInt(0)

		fmt.Println("nonce:  ", statedb.GetNonce(callerAddress))
		fmt.Println("balance:  ", statedb.GetBalance(callerAddress))
		_, contractAddress, _, err := evm.Create(vm.AccountRef(callerAddress), contractCodeBytes, 5000000, value)
		check(err)
		fmt.Println("contractAddress:  ", contractAddress)

		for i := 0; i < numRuns; i++ {
			snapshot := statedb.Snapshot()
			start := time.Now()

			_, _, err := evm.Call(vm.AccountRef(callerAddress), contractAddress, calldataBytes, gasLimit, value, false)

			timeTaken := time.Since(start)

			fmt.Println("result:   ", float64(timeTaken.Microseconds())/1e3)
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
