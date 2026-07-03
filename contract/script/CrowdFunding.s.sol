// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Script} from "forge-std/Script.sol";
import {CrowdFunding} from "../src/CrowdFunding.sol";

contract CrowdFundingScript is Script {
    function run() external returns (CrowdFunding) {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");

        vm.startBroadcast(deployerPrivateKey);
        CrowdFunding crowdFunding = new CrowdFunding();
        vm.stopBroadcast();

        return crowdFunding;
    }
}
