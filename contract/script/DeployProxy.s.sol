// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Script} from "forge-std/Script.sol";
import {CrowdFunding} from "../src/CrowdFunding.sol";
import {CrowdFundingProxy} from "../src/CrowdFundingProxy.sol";

contract DeployProxyScript is Script {
    function run()
        external
        returns (address proxyAddress, address implementationAddress)
    {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address deployer = vm.addr(deployerPrivateKey);

        vm.startBroadcast(deployerPrivateKey);

        CrowdFunding implementation = new CrowdFunding();
        CrowdFundingProxy proxy = new CrowdFundingProxy(
            address(implementation),
            deployer
        );

        vm.stopBroadcast();

        return (address(proxy), address(implementation));
    }
}
