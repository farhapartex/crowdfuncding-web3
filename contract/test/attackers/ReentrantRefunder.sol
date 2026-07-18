// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {CrowdFunding} from "../../src/CrowdFunding.sol";

contract ReentrantRefunder {
    CrowdFunding public immutable crowdFunding;
    uint256 public campaignId;
    bool public hasAttacked;

    constructor(CrowdFunding _crowdFunding) {
        crowdFunding = _crowdFunding;
    }

    function contribute(uint256 _campaignId) external payable {
        campaignId = _campaignId;
        crowdFunding.contributeEth{value: msg.value}(_campaignId);
    }

    function attack() external {
        crowdFunding.refund(campaignId);
    }

    receive() external payable {
        if (!hasAttacked) {
            hasAttacked = true;
            crowdFunding.refund(campaignId);
        }
    }
}
