// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {CrowdFunding} from "../../src/CrowdFunding.sol";

contract CrowdFundingV2 is CrowdFunding {
    function totalRaised() external view returns (uint256 total) {
        uint256 count = campaigns.length;
        for (uint256 i = 0; i < count; i++) {
            total += campaigns[i].amountRaised;
        }
    }
}
