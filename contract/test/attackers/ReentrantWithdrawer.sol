// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {CrowdFunding} from "../../src/CrowdFunding.sol";

contract ReentrantWithdrawer {
    CrowdFunding public immutable crowdFunding;
    uint256 public campaignId;
    bool public hasAttacked;

    constructor(CrowdFunding _crowdFunding) {
        crowdFunding = _crowdFunding;
    }

    function createCampaign(string calldata title, string calldata description, uint256 goal, uint256 duration)
        external
        returns (uint256)
    {
        campaignId = crowdFunding.createCampaign(title, description, goal, duration);
        return campaignId;
    }

    function attack() external {
        crowdFunding.withdraw(campaignId);
    }

    receive() external payable {
        if (!hasAttacked) {
            hasAttacked = true;
            crowdFunding.withdraw(campaignId);
        }
    }
}
