package main

// import (
// "fmt"
// "math/big"
// "math/rand"
// "os"
// "time"

// "github.com/holiman/uint256"
// libcommon "github.com/ledgerwatch/erigon-lib/common"
// "github.com/ledgerwatch/erigon-lib/common/hexutility"
// types2 "github.com/ledgerwatch/erigon-lib/types"
// "github.com/ledgerwatch/erigon/common"
// "github.com/ledgerwatch/erigon/core"
// "github.com/ledgerwatch/erigon/core/types"
// "github.com/ledgerwatch/erigon/core/vm"
// "github.com/ledgerwatch/erigon/core/vm/evmtypes"
// "github.com/ledgerwatch/erigon/params"
// "github.com/spf13/cobra"
// )

// var (
// 	contractCodePath string
// 	calldata         string
// 	numRuns          int
// )

// func check(e error) {
// 	if e != nil {
// 		fmt.Fprintln(os.Stderr, e)
// 		os.Exit(1)
// 	}
// }

// func randIntInRange(min, max int) int {
// 	return (rand.Intn(max-min) + min)
// }

// func randHash() libcommon.Hash {
// 	var h libcommon.Hash
// 	for i := 0; i < 32; i++ {
// 		h[i] = byte(rand.Intn(255))
// 	}
// 	return h
// }

// func randAddr() *libcommon.Address {
// 	var a libcommon.Address
// 	for j := 0; j < 20; j++ {
// 		a[j] = byte(rand.Intn(255))
// 	}
// 	return &a
// }

// func randAccessList() types2.AccessList {
// 	size := randIntInRange(4, 10)
// 	var result types2.AccessList
// 	for i := 0; i < size; i++ {
// 		var tup types2.AccessTuple

// 		tup.Address = *randAddr()
// 		tup.StorageKeys = append(tup.StorageKeys, randHash())
// 		result = append(result, tup)
// 	}
// 	return result
// }

// var cmd = &cobra.Command{
// 	Use:   "runner-erigon",
// 	Short: "erigon runner for evm-bench",
// 	Run: func(_ *cobra.Command, _ []string) {
// 		contractCodeHex, err := os.ReadFile(contractCodePath)
// 		check(err)

// 		contractCodeBytes := hexutility.Hex2Bytes(string(contractCodeHex))
// 		calldataBytes := hexutility.Hex2Bytes(calldata)

// 		// fmt.Println(contractCodeBytes)
// 		// fmt.Println(calldataBytes)

// 		zeroAddress := libcommon.BytesToAddress(common.FromHex("0x0000000000000000000000000000000000000000"))
// 		callerAddress := libcommon.BytesToAddress(common.FromHex("0x1000000000000000000000000000000000000001"))
// 		// callerAddress1 := libcommon.BytesToAddress(common.FromHex("0x1000000000000000000000000000000000000002"))
// 		// fmt.Println(zeroAddress)
// 		// fmt.Println(callerAddress)

// 		var config = params.MainnetChainConfig
// 		rules := config.Rules(uint64(12965000), uint64(1681338458))
// 		// fmt.Println(config)
// 		// fmt.Println(rules)
// 		gasLimit := ^uint64(1000000000)

// 		//create genesis
// 		address1 := libcommon.HexToAddress("0x000d836201318ec6899a67540690382780743280")
// 		address2 := libcommon.HexToAddress("0x001762430ea9c3a26e5749afdb70da5f78ddbb8c")
// 		funds := new(big.Int)
// 		funds.SetString("ad78ebc5ac6200000", 16)

// 		alloc := types.GenesisAlloc{
// 			address1: {Balance: funds},
// 			address2: {Balance: funds},
// 			//callerAddress: {Balance: funds},
// 		}

// 		genesis := &types.Genesis{
// 			Config:     config,
// 			Coinbase:   zeroAddress,
// 			Difficulty: big.NewInt(17179869184),
// 			GasLimit:   5000000000,
// 			Number:     1681338457,
// 			Timestamp:  1681338458,
// 			Alloc:      alloc,
// 		}
// 		// fmt.Println(genesis)

// 		block, statedb, err := core.GenesisToBlock(genesis, "")
// 		// testkv := memdb.New("")
// 		// tx, _ := testkv.BeginRw(context.Background())
// 		// r := state.NewPlainState(tx, 0, nil)
// 		// statedb := state.New(r)

// 		check(err)

// 		var accesslist types2.AccessList
// 		accesslist = randAccessList()
// 		statedb.Prepare(rules, callerAddress, zeroAddress, &zeroAddress, vm.ActivePrecompiles(rules), accesslist)
// 		blockContext := core.NewEVMBlockContext(block.Header(), core.GetHashFn(block.Header(), nil), nil, &zeroAddress)

// 		fmt.Println("!!", blockContext)

// 		var value *uint256.Int
// 		value = uint256.NewInt(0)
// 		evm := vm.NewEVM(blockContext, evmtypes.TxContext{}, statedb, config, vm.Config{})

// 		//evm.intraBlockState.GetNonce(caller.Address())
// 		fmt.Println("nonce:  ", statedb.GetNonce(address2))
// 		_, contractAddress, _, err := evm.Create(vm.AccountRef(address1), contractCodeBytes, gasLimit, value)

// 		check(err)
// 		fmt.Println("??", contractAddress)

// 		for i := 0; i < numRuns; i++ {
// 			snapshot := statedb.Snapshot()
// 			// statedb.AddAddressToAccessList(msg.From)
// 			// statedb.AddAddressToAccessList(*msg.To)

// 			start := time.Now()

// 			_, _, err := evm.Call(vm.AccountRef(callerAddress), contractAddress, calldataBytes, gasLimit, value, false)

// 			timeTaken := time.Since(start)

// 			fmt.Println(float64(timeTaken.Microseconds()) / 1e3)

// 			check(err)

// 			statedb.RevertToSnapshot(snapshot)
// 		}

// 	},
// }

// func init() {
// 	cmd.Flags().StringVar(&contractCodePath, "contract-code-path", "", "Path to the hex contract code to deploy and run")
// 	cmd.MarkFlagRequired("contract-code-path")
// 	cmd.Flags().StringVar(&calldata, "calldata", "", "Hex of calldata to use when calling the contract")
// 	cmd.MarkFlagRequired("calldata")
// 	cmd.Flags().IntVar(&numRuns, "num-runs", 0, "Number of times to run the benchmark")
// 	cmd.MarkFlagRequired("num-runs")
// }

func bakup() {
	// if err := cmd.Execute(); err != nil {
	// 	fmt.Fprintln(os.Stderr, err)
	// 	os.Exit(1)
	// }
}
