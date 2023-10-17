// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.5.2;

contract TenThousandHashes {
    function Benchmark() external pure {
        for (uint256 i = 0; i < 20000; i++) {
            keccak256(abi.encodePacked(i));
        }
    }
}
